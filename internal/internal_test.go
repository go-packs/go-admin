package internal

import (
	"testing"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func setupTestDB() (*gorm.DB, *admin.Registry) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.AdminUser{}, &models.Permission{}, &models.AuditLog{}, &MockModel{}); err != nil {
		panic(err)
	}
	reg := admin.NewRegistry(db)
	return db, reg
}

func TestAuthLogic(t *testing.T) {
	_, reg := setupTestDB()

	t.Run("IsAllowed", func(t *testing.T) {
		if !IsAllowed(reg, "admin", "User", "list") {
			t.Error("Admin should be allowed everything")
		}
		reg.DB.Create(&models.Permission{Role: "editor", ResourceName: "User", Action: "list"})
		if !IsAllowed(reg, "editor", "User", "list") {
			t.Error("Editor should be allowed list on User")
		}
		if IsAllowed(reg, "editor", "User", "delete") {
			t.Error("Editor should not be allowed delete on User")
		}
	})
}

func TestAuditLogic(t *testing.T) {
	db, reg := setupTestDB()
	user := &models.AdminUser{Email: "admin@example.com"}

	t.Run("RecordAction", func(t *testing.T) {
		RecordAction(reg, user, "User", "1", "Update", "Changed email")
		var log models.AuditLog
		if err := db.First(&log).Error; err != nil {
			t.Fatalf("Failed to find audit log: %v", err)
		}
		if log.UserEmail != user.Email || log.Action != "Update" {
			t.Error("Audit log record mismatch")
		}
	})
}

func TestCRUDLogic(t *testing.T) {
	_, reg := setupTestDB()
	reg.Register(MockModel{})

	t.Run("Lifecycle", func(t *testing.T) {
		item := &MockModel{Name: "Initial"}
		if err := Create(reg, item); err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		fetched, _ := Get(reg, "MockModel", item.ID)
		if fetched.(*MockModel).Name != "Initial" {
			t.Error("Get failed")
		}

		item.Name = "Updated"
		if err := Update(reg, item); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		list, _ := List(reg, "MockModel")
		if len(list.([]MockModel)) != 1 {
			t.Error("List failed")
		}

		if err := Delete(reg, "MockModel", item.ID); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		listAfter, _ := List(reg, "MockModel")
		if len(listAfter.([]MockModel)) != 0 {
			t.Error("Delete did not remove item")
		}
	})
}
