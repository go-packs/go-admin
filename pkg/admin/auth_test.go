package admin

import (
	"testing"
	"github.com/ajeet-kumar1087/go-admin/pkg/admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPasswordHashing(t *testing.T) {
	user := &models.AdminUser{}
	password := "secure123"
	
	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}

	if user.PasswordHash == password {
		t.Errorf("Password hash should not be plain text")
	}

	if !user.CheckPassword(password) {
		t.Errorf("CheckPassword failed for correct password")
	}

	if user.CheckPassword("wrongpassword") {
		t.Errorf("CheckPassword should fail for incorrect password")
	}
}

func TestIsAllowed(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Permission{})
	reg := NewRegistry(db)

	// Admin role bypass
	if !reg.IsAllowed("admin", "AnyResource", "anyAction") {
		t.Errorf("Admin role should always be allowed")
	}

	// Guest role (no permissions)
	if reg.IsAllowed("guest", "User", "list") {
		t.Errorf("Guest should not have list permission by default")
	}

	// Seed specific permission
	db.Create(&models.Permission{Role: "editor", ResourceName: "Product", Action: "edit"})

	if !reg.IsAllowed("editor", "Product", "edit") {
		t.Errorf("Editor should have edit permission for Product")
	}

	if reg.IsAllowed("editor", "Product", "delete") {
		t.Errorf("Editor should not have delete permission for Product")
	}
}
