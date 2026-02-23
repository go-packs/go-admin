package admin_test

import (
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
