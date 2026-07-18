package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	admin "github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HandlerProduct struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Price  float64
	Stock  int
	Active bool
}

func setupTestDB() (*gorm.DB, *admin.Registry) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.AdminUser{}, &models.Session{}, &models.Permission{}, &models.AuditLog{}, &HandlerProduct{}); err != nil {
		panic(err)
	}
	reg := admin.NewRegistry(db)
	return db, reg
}

func TestAuthHandlers(t *testing.T) {
	db, reg := setupTestDB()

	t.Run("LoginSuccess", func(t *testing.T) {
		user := &models.AdminUser{Email: "test@example.com"}
		if err := user.SetPassword("password123"); err != nil {
			t.Fatalf("set password: %v", err)
		}
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user: %v", err)
		}

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

func TestResourceHandlers(t *testing.T) {
	db, reg := setupTestDB()
	res := reg.Register(HandlerProduct{}).
		RegisterField("ID", "ID", true).
		RegisterField("Name", "Name", false).
		RegisterField("Price", "Price", false).
		RegisterField("Stock", "Stock", false).
		RegisterField("Active", "Active", false)
	user := &models.AdminUser{Email: "admin@example.com", Role: "admin"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	t.Run("HandleSaveCreatesTypedModelAndAuditLog", func(t *testing.T) {
		form := url.Values{}
		form.Set("Name", "Keyboard")
		form.Set("Price", "149.95")
		form.Set("Stock", "7")
		form.Set("Active", "on")
		req := httptest.NewRequest("POST", "/admin/HandlerProduct/save", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		HandleSave(reg, res, w, req, user)

		if w.Code != http.StatusSeeOther {
			t.Fatalf("Expected redirect, got %d", w.Code)
		}
		var product HandlerProduct
		if err := db.First(&product).Error; err != nil {
			t.Fatalf("find product: %v", err)
		}
		if product.Name != "Keyboard" || product.Price != 149.95 || product.Stock != 7 || !product.Active {
			t.Fatalf("Product fields were not saved from form: %#v", product)
		}
		var log models.AuditLog
		if err := db.Where("resource_name = ? AND action = ?", "HandlerProduct", "Create").First(&log).Error; err != nil {
			t.Fatalf("find audit log: %v", err)
		}
		if log.UserEmail != user.Email || log.RecordID == "" {
			t.Fatalf("Unexpected audit log: %#v", log)
		}
	})

	t.Run("HandleSaveUpdatesExistingModel", func(t *testing.T) {
		product := HandlerProduct{Name: "Mouse", Price: 20, Stock: 2}
		if err := db.Create(&product).Error; err != nil {
			t.Fatalf("create product: %v", err)
		}
		form := url.Values{}
		form.Set("ID", fmt.Sprintf("%d", product.ID))
		form.Set("Name", "Gaming Mouse")
		form.Set("Price", "35.5")
		form.Set("Stock", "4")
		req := httptest.NewRequest("POST", "/admin/HandlerProduct/save", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		HandleSave(reg, res, w, req, user)

		var updated HandlerProduct
		if err := db.First(&updated, product.ID).Error; err != nil {
			t.Fatalf("find updated product: %v", err)
		}
		if updated.Name != "Gaming Mouse" || updated.Price != 35.5 || updated.Stock != 4 {
			t.Fatalf("Product was not updated from form: %#v", updated)
		}
		var count int64
		db.Model(&models.AuditLog{}).Where("resource_name = ? AND action = ?", "HandlerProduct", "Update").Count(&count)
		if count == 0 {
			t.Fatal("Expected update audit log")
		}
	})

	t.Run("HandleSearchAPIUsesRegisteredTextFields", func(t *testing.T) {
		if err := db.Create(&HandlerProduct{Name: "Standing Desk", Price: 300}).Error; err != nil {
			t.Fatalf("create search product: %v", err)
		}
		req := httptest.NewRequest("GET", "/admin/HandlerProduct/search?q=Desk", nil)
		w := httptest.NewRecorder()

		HandleSearchAPI(reg, "HandlerProduct", w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var results []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
			t.Fatalf("decode search results: %v", err)
		}
		if len(results) == 0 || results[0]["text"] != "Standing Desk" {
			t.Fatalf("Unexpected search results: %#v", results)
		}
	})

	t.Run("HandleExportWritesCSV", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/HandlerProduct/export", nil)
		w := httptest.NewRecorder()

		HandleExport(reg, res, w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		rows, err := csv.NewReader(strings.NewReader(w.Body.String())).ReadAll()
		if err != nil {
			t.Fatalf("read csv: %v", err)
		}
		if len(rows) < 2 {
			t.Fatalf("Expected header and rows, got %#v", rows)
		}
		if rows[0][1] != "Name" {
			t.Fatalf("Unexpected CSV header: %#v", rows[0])
		}
	})
}
