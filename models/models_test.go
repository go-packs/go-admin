package models

import (
	"testing"
)

func TestAdminUser(t *testing.T) {
	t.Run("PasswordHashing", func(t *testing.T) {
		user := &AdminUser{Email: "test@example.com"}
		err := user.SetPassword("password123")
		if err != nil {
			t.Fatal(err)
		}
		if user.PasswordHash == "" || user.PasswordHash == "password123" {
			t.Errorf("Password was not hashed correctly")
		}
		if !user.CheckPassword("password123") {
			t.Errorf("Password check failed for correct password")
		}
		if user.CheckPassword("wrongpassword") {
			t.Errorf("Password check passed for incorrect password")
		}
	})
}
