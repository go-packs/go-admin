package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, *admin.Registry) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.AdminUser{}, &models.Session{}, &models.Permission{}); err != nil {
		panic(err)
	}
	reg := admin.NewRegistry(db)
	return db, reg
}

func TestRouter(t *testing.T) {
	_, reg := setupTestDB()
	router := NewRouter(reg)

	t.Run("LoginPath", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("AuthGuardRedirect", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected 303, got %d", w.Code)
		}
		if w.Header().Get("Location") != "/admin/login" {
			t.Errorf("Expected redirect to login, got %s", w.Header().Get("Location"))
		}
	})
}
