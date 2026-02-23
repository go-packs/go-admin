package handlers

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupTestDB() (*gorm.DB, *admin.Registry) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.AdminUser{}, &models.Session{}, &models.Permission{})
	reg := admin.NewRegistry(db)
	return db, reg
}

func TestAuthHandlers(t *testing.T) {
	db, reg := setupTestDB()

	t.Run("LoginSuccess", func(t *testing.T) {
		user := &models.AdminUser{Email: "test@example.com"}
		user.SetPassword("password123")
		db.Create(user)

		data := url.Values{}
		data.Set("email", "test@example.com")
		data.Set("password", "password123")

		req := httptest.NewRequest("POST", "/admin/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		Login(reg)(w, req)

		if w.Code != 303 {
			t.Errorf("Expected 303, got %d", w.Code)
		}
		if !strings.Contains(w.Header().Get("Set-Cookie"), "admin_session") {
			t.Error("Session cookie not set")
		}
	})

	t.Run("Logout", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/logout", nil)
		cookie := &http.Cookie{Name: "admin_session", Value: "test-sess", Path: "/admin"}
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		Logout(reg)(w, req)

		if w.Code != 303 {
			t.Errorf("Expected 303, got %d", w.Code)
		}
		if !strings.Contains(w.Header().Get("Set-Cookie"), "Max-Age=0") {
			t.Error("Session cookie not expired")
		}
	})
}
