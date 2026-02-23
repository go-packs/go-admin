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
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&TestModel{}, &admin.Permission{})
	reg := admin.NewRegistry(db)

	t.Run("RegistryInitialization", func(t *testing.T) {
		if reg.DB != db { t.Error("DB not set") }
		res := reg.Register(TestModel{})
		if res.Name != "TestModel" { t.Errorf("Expected TestModel, got %s", res.Name) }
	})

	t.Run("Authentication", func(t *testing.T) {
		user := &admin.AdminUser{}
		user.SetPassword("123")
		if !user.CheckPassword("123") { t.Error("Password check failed") }
		
		if !internal.IsAllowed(reg, "admin", "Any", "Any") { t.Error("Admin should be allowed") }
		db.Create(&admin.Permission{Role: "editor", ResourceName: "Product", Action: "edit"})
		if !internal.IsAllowed(reg, "editor", "Product", "edit") { t.Error("Permission check failed") }
	})

	t.Run("CRUD", func(t *testing.T) {
		reg.Register(TestModel{})
		item := &TestModel{Name: "Go"}
		internal.Create(reg, item)
		
		fetched, _ := internal.Get(reg, "TestModel", item.ID)
		if fetched.(*TestModel).Name != "Go" { t.Error("Create/Get failed") }
		
		item.Name = "Rust"
		internal.Update(reg, item)
		fetched, _ = internal.Get(reg, "TestModel", item.ID)
		if fetched.(*TestModel).Name != "Rust" { t.Error("Update failed") }
		
		internal.Delete(reg, "TestModel", item.ID)
		list, _ := internal.List(reg, "TestModel")
		if len(list.([]TestModel)) != 0 { t.Error("Delete failed") }
	})

	t.Run("Configuration", func(t *testing.T) {
		yaml := "site_title: 'Custom'\ndefault_per_page: 50"
		os.WriteFile("test.yml", []byte(yaml), 0644)
		defer os.Remove("test.yml")
		
		conf, _ := admin.LoadConfig("test.yml")
		if conf.SiteTitle != "Custom" || conf.DefaultPerPage != 50 { t.Error("Config load failed") }
	})
}
