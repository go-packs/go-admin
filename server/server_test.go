package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RouterProduct struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func setupTestDB() (*gorm.DB, *admin.Registry) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.AdminUser{}, &models.Session{}, &models.Permission{}, &RouterProduct{}); err != nil {
		panic(err)
	}
	reg := admin.NewRegistry(db)
	return db, reg
}

func createSession(t *testing.T, db *gorm.DB, email, role string) *http.Cookie {
	t.Helper()
	user := &models.AdminUser{Email: email, Role: role}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	session := &models.Session{ID: email + "-session", UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour)}
	if err := db.Create(session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}
	return &http.Cookie{Name: "admin_session", Value: session.ID, Path: "/admin"}
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

	t.Run("AuthenticatedDashboard", func(t *testing.T) {
		db, reg := setupTestDB()
		router := NewRouter(reg)
		req := httptest.NewRequest("GET", "/admin/", nil)
		req.AddCookie(createSession(t, db, "admin@example.com", "admin"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("AuthenticatedResourceList", func(t *testing.T) {
		db, reg := setupTestDB()
		reg.Register(RouterProduct{}).
			RegisterField("ID", "ID", true).
			RegisterField("Name", "Name", false)
		if err := db.Create(&RouterProduct{Name: "Keyboard"}).Error; err != nil {
			t.Fatalf("create product: %v", err)
		}
		router := NewRouter(reg)
		req := httptest.NewRequest("GET", "/admin/RouterProduct", nil)
		req.AddCookie(createSession(t, db, "admin-list@example.com", "admin"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("ForbiddenWithoutPermission", func(t *testing.T) {
		db, reg := setupTestDB()
		reg.Register(RouterProduct{}).RegisterField("ID", "ID", true)
		router := NewRouter(reg)
		req := httptest.NewRequest("GET", "/admin/RouterProduct", nil)
		req.AddCookie(createSession(t, db, "viewer@example.com", "viewer"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected 403, got %d", w.Code)
		}
	})

	t.Run("UnknownResourceReturnsNotFound", func(t *testing.T) {
		db, reg := setupTestDB()
		router := NewRouter(reg)
		req := httptest.NewRequest("GET", "/admin/Missing", nil)
		req.AddCookie(createSession(t, db, "missing@example.com", "admin"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected 404, got %d", w.Code)
		}
	})
}
