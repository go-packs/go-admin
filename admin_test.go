package admin_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/internal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func TestCore(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&TestModel{}, &admin.Permission{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	reg := admin.NewRegistry(db)

	t.Run("RegistryInitialization", func(t *testing.T) {
		if reg.DB != db {
			t.Error("DB not set")
		}
		res := reg.Register(TestModel{})
		if res.Name != "TestModel" {
			t.Errorf("Expected TestModel, got %s", res.Name)
		}
	})

	t.Run("Authentication", func(t *testing.T) {
		user := &admin.AdminUser{}
		if err := user.SetPassword("123"); err != nil {
			t.Fatalf("set password: %v", err)
		}
		if !user.CheckPassword("123") {
			t.Error("Password check failed")
		}

		if !internal.IsAllowed(reg, "admin", "Any", "Any") {
			t.Error("Admin should be allowed")
		}
		if err := db.Create(&admin.Permission{Role: "editor", ResourceName: "Product", Action: "edit"}).Error; err != nil {
			t.Fatalf("create permission: %v", err)
		}
		if !internal.IsAllowed(reg, "editor", "Product", "edit") {
			t.Error("Permission check failed")
		}
	})

	t.Run("CRUD", func(t *testing.T) {
		reg.Register(TestModel{})
		item := &TestModel{Name: "Go"}
		if err := internal.Create(reg, item); err != nil {
			t.Fatalf("create item: %v", err)
		}

		fetched, err := internal.Get(reg, "TestModel", item.ID)
		if err != nil {
			t.Fatalf("get item: %v", err)
		}
		if fetched.(*TestModel).Name != "Go" {
			t.Error("Create/Get failed")
		}

		item.Name = "Rust"
		if err := internal.Update(reg, item); err != nil {
			t.Fatalf("update item: %v", err)
		}
		fetched, err = internal.Get(reg, "TestModel", item.ID)
		if err != nil {
			t.Fatalf("get item: %v", err)
		}
		if fetched.(*TestModel).Name != "Rust" {
			t.Error("Update failed")
		}

		if err := internal.Delete(reg, "TestModel", item.ID); err != nil {
			t.Fatalf("delete item: %v", err)
		}
		list, err := internal.List(reg, "TestModel")
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}
		if len(list.([]TestModel)) != 0 {
			t.Error("Delete failed")
		}
	})

	t.Run("Configuration", func(t *testing.T) {
		yaml := "site_title: 'Custom'\ndefault_per_page: 50"
		if err := os.WriteFile("test.yml", []byte(yaml), 0644); err != nil {
			t.Fatalf("write conf: %v", err)
		}
		defer func() {
			_ = os.Remove("test.yml")
		}()

		conf, err := admin.LoadConfig("test.yml")
		if err != nil {
			t.Fatalf("load conf: %v", err)
		}
		if conf.SiteTitle != "Custom" || conf.DefaultPerPage != 50 {
			t.Error("Config load failed")
		}
	})
}

func TestRegistryHelpers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	reg := admin.NewRegistry(db)

	t.Run("GroupsResourcesAndPages", func(t *testing.T) {
		reg.Register(TestModel{}).SetGroup("Content")
		reg.AddPage("Status", "Operations", func(w http.ResponseWriter, r *http.Request) {})

		resourceGroups := reg.GetGroupedResources()
		if len(resourceGroups["Content"]) != 1 || resourceGroups["Content"][0].Name != "TestModel" {
			t.Fatalf("Unexpected resource groups: %#v", resourceGroups)
		}

		pageGroups := reg.GetGroupedPages()
		if len(pageGroups["Operations"]) != 1 || pageGroups["Operations"][0].Name != "Status" {
			t.Fatalf("Unexpected page groups: %#v", pageGroups)
		}
	})

	t.Run("ResourceNamesAndCharts", func(t *testing.T) {
		names := reg.ResourceNames()
		if len(names) != 1 || names[0] != "TestModel" {
			t.Fatalf("Unexpected resource names: %#v", names)
		}

		reg.AddChart("Records", "bar", func(db *gorm.DB) ([]string, []float64) {
			return []string{"A"}, []float64{1}
		})
		if len(reg.Charts) != 1 || reg.Charts[0].Label != "Records" || reg.Charts[0].Type != "bar" {
			t.Fatalf("Unexpected chart registration: %#v", reg.Charts)
		}
	})

	t.Run("FlashCookieRoundTrip", func(t *testing.T) {
		w := httptest.NewRecorder()
		reg.SetFlash(w, "Saved")
		cookies := w.Result().Cookies()
		if len(cookies) != 1 || cookies[0].Name != "admin_flash" || cookies[0].Value != "Saved" {
			t.Fatalf("Unexpected flash cookie: %#v", cookies)
		}

		req := httptest.NewRequest("GET", "/admin", nil)
		req.AddCookie(cookies[0])
		clearWriter := httptest.NewRecorder()
		if got := reg.GetFlash(clearWriter, req); got != "Saved" {
			t.Fatalf("Expected flash Saved, got %q", got)
		}
		clearCookies := clearWriter.Result().Cookies()
		if len(clearCookies) != 1 || clearCookies[0].MaxAge != -1 {
			t.Fatalf("Expected flash clearing cookie, got %#v", clearCookies)
		}
	})
}
