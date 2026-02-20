package main

import (
	"fmt"
	"github.com/ajeet-kumar1087/go-admin/pkg/admin"
	"github.com/ajeet-kumar1087/go-admin/pkg/admin/config"
	"github.com/ajeet-kumar1087/go-admin/pkg/admin/models"
	"github.com/ajeet-kumar1087/go-admin/pkg/admin/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
)

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Email string
	Role  string
}

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string
	Price float64
	Image string
}

type ProductInfo struct {
	ID          uint `gorm:"primaryKey"`
	ProductID   uint
	Description string
	Manufacturer string
}

func main() {
	db, err := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate schemas
	db.AutoMigrate(&User{}, &Product{}, &ProductInfo{}, &models.Permission{}, &Role{}, &models.AdminUser{}, &models.Session{}, &models.AuditLog{})

	adm := admin.NewRegistry(db)

	// Load Custom Config
	conf, err := config.LoadConfig("admin.yml")
	if err == nil {
		adm.SetConfig(conf)
	}

	roles := []string{"admin", "editor", "viewer"}

	// Activity Action Helper
	addActivityAction := func(r *resource.Resource) {
		r.AddMemberAction("activity", "View History", func(res *resource.Resource, w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			http.Redirect(w, r, fmt.Sprintf("/admin/AuditLog?q_ResourceName=%s&q_RecordID=%s", res.Name, id), http.StatusSeeOther)
		})
	}

	// Administration Group
	adm.Register(models.AdminUser{}).
		SetGroup("Administration").
		RegisterField("ID", "ID", true).
		RegisterField("Email", "Email", false).
		RegisterField("Role", "Role", false).
		SetFieldType("Role", "select", roles...)

	adm.Register(models.AuditLog{}).
		SetGroup("Administration").
		RegisterField("ID", "ID", true).
		RegisterField("CreatedAt", "Time", true).
		RegisterField("UserEmail", "User", true).
		RegisterField("ResourceName", "Resource", true).
		RegisterField("RecordID", "Record ID", true).
		RegisterField("Action", "Action", true).
		RegisterField("Changes", "Changes", true)

	adm.Register(Role{}).
		SetGroup("Administration").
		RegisterField("ID", "ID", true).
		RegisterField("Name", "Role Name", false)

	adm.Register(models.Permission{}).
		SetGroup("Administration").
		RegisterField("ID", "ID", true).
		RegisterField("Role", "Role Name", false).
		RegisterField("ResourceName", "Resource", false).
		RegisterField("Action", "Action", false).
		SetFieldType("Role", "select", roles...).
		SetFieldType("ResourceName", "select", adm.ResourceNames()...).
		SetFieldType("Action", "select", "list", "show", "new", "edit", "save", "delete")

	// Users (Default Group)
	uRes := adm.Register(User{}).
		RegisterField("ID", "ID", true).
		RegisterField("Email", "Email Address", false).
		RegisterField("Role", "User Role", false).
		SetFieldType("Role", "select", roles...).
		SetDecorator("Role", func(val interface{}) template.HTML {
			role := val.(string)
			color := "#64748b" // default
			if role == "admin" { color = "#ef4444" } else if role == "editor" { color = "#3b82f6" }
			return template.HTML(fmt.Sprintf(`<span style="background: %s; color: white; padding: 0.2rem 0.5rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 600;">%s</span>`, color, role))
		}).
		AddMemberAction("activate", "Activate User", func(res *resource.Resource, w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			fmt.Fprintf(w, "User %s has been activated! (Simulated)", id)
		}).
		AddScope("admins", "Admins Only", func(db *gorm.DB) *gorm.DB {
			return db.Where("role = ?", "admin")
		}).
		AddScope("editors", "Editors Only", func(db *gorm.DB) *gorm.DB {
			return db.Where("role = ?", "editor")
		})
	addActivityAction(uRes)

	// Products Group
	pRes := adm.Register(Product{}).
		SetGroup("Products").
		RegisterField("ID", "ID", true).
		RegisterField("Name", "Product Name", false).
		RegisterField("Price", "Price", false).
		RegisterField("Image", "Product Image", false).
		SetFieldType("Price", "number").
		SetFieldType("Image", "image").
		SetDecorator("Price", func(val interface{}) template.HTML {
			price := val.(float64)
			return template.HTML(fmt.Sprintf("<strong>$%.2f</strong>", price))
		}).
		HasMany("ProductInfo", "Technical Specifications", "ProductInfo", "ProductID").
		AddCollectionAction("discount", "Apply 10% Bulk Discount", func(res *resource.Resource, w http.ResponseWriter, r *http.Request) {
			db.Model(&Product{}).Where("price > ?", 0).Update("price", gorm.Expr("price * 0.9"))
			http.Redirect(w, r, "/admin/Product", http.StatusSeeOther)
		}).
		AddBatchAction("batch_delete", "Delete Selected", func(res *resource.Resource, ids []string, w http.ResponseWriter, r *http.Request) {
			db.Where("id IN ?", ids).Delete(&Product{})
			http.Redirect(w, r, "/admin/Product", http.StatusSeeOther)
		}).
		AddScope("expensive", "Expensive (> $500)", func(db *gorm.DB) *gorm.DB {
			return db.Where("price > ?", 500)
		}).
		AddScope("cheap", "Cheap (< $100)", func(db *gorm.DB) *gorm.DB {
			return db.Where("price < ?", 100)
		})
	addActivityAction(pRes)

	piRes := adm.Register(ProductInfo{}).
		SetGroup("Products").
		RegisterField("ID", "ID", true).
		RegisterField("ProductID", "Product", false).
		RegisterField("Description", "Description", false).
		RegisterField("Manufacturer", "Manufacturer", false).
		BelongsTo("ProductID", "Parent Product", "Product", "ID").
		SetSearchable("ProductID", "Product")
	addActivityAction(piRes)

	// 4. Register Dashboard Charts
	adm.AddChart("Users by Role", "pie", func(db *gorm.DB) ([]string, []float64) {
		var results []struct {
			Role  string
			Count int64
		}
		db.Model(&User{}).Select("role, count(*) as count").Group("role").Scan(&results)
		labels := []string{}; values := []float64{}
		for _, r := range results { labels = append(labels, r.Role); values = append(values, float64(r.Count)) }
		return labels, values
	})

	adm.AddChart("Price Distribution", "bar", func(db *gorm.DB) ([]string, []float64) {
		var results []struct {
			Name  string
			Price float64
		}
		db.Model(&Product{}).Limit(5).Scan(&results)
		labels := []string{}; values := []float64{}
		for _, r := range results { labels = append(labels, r.Name); values = append(values, r.Price) }
		return labels, values
	})

	// 5. Register Custom Pages
	adm.AddPage("SystemStatus", "Administration", func(w http.ResponseWriter, r *http.Request) {
		content := template.HTML(`
			<div style="background: white; border-radius: 0.5rem; overflow: hidden;">
				<table style="width: 100%; border-collapse: collapse;">
					<tr style="border-bottom: 1px solid #e2e8f0;"><td style="padding: 1rem; font-weight: 600;">Server Uptime</td><td style="padding: 1rem;">12 days, 4 hours</td></tr>
					<tr style="border-bottom: 1px solid #e2e8f0;"><td style="padding: 1rem; font-weight: 600;">Database Status</td><td style="padding: 1rem; color: #10b981;">Connected (v15.2)</td></tr>
					<tr style="border-bottom: 1px solid #e2e8f0;"><td style="padding: 1rem; font-weight: 600;">Storage Used</td><td style="padding: 1rem;">45.2 GB / 100 GB</td></tr>
					<tr><td style="padding: 1rem; font-weight: 600;">Go Version</td><td style="padding: 1rem;">1.25.7</td></tr>
				</table>
			</div>
		`)
		adm.RenderCustomPage(w, r, "System Status Overview", content)
	})

	// 6. Seed Data
	var adminCount int64; db.Model(&models.AdminUser{}).Count(&adminCount)
	if adminCount == 0 {
		adminUser := &models.AdminUser{Email: "admin@example.com", Role: "admin"}
		adminUser.SetPassword("password123"); db.Create(adminUser)
		db.Create(&Role{Name: "admin"}); db.Create(&Role{Name: "editor"}); db.Create(&Role{Name: "viewer"})
		db.Create(&models.Permission{Role: "editor", ResourceName: "Product", Action: "list"})
		db.Create(&models.Permission{Role: "editor", ResourceName: "Product", Action: "show"})
		db.Create(&models.Permission{Role: "editor", ResourceName: "Product", Action: "edit"})
		db.Create(&models.Permission{Role: "editor", ResourceName: "Product", Action: "save"})
		p1 := &Product{Name: "Mechanical Keyboard", Price: 150.00}
		p2 := &Product{Name: "Gaming Mouse", Price: 80.00}
		db.Create(p1); db.Create(p2)
		db.Create(&ProductInfo{ProductID: p1.ID, Description: "Blue Switches", Manufacturer: "Razer"})
		db.Create(&ProductInfo{ProductID: p1.ID, Description: "RGB Lighting", Manufacturer: "Razer"})
		db.Create(&User{Email: "user@example.com", Role: "editor"})
	}

	port := ":8080"
	fmt.Printf("\n🚀 Admin panel running at http://localhost%s/admin\n", port)
	http.Handle("/admin/", adm)
	log.Fatal(http.ListenAndServe(port, nil))
}
