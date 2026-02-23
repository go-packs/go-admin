package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"github.com/go-packs/go-admin/view"
	"github.com/google/uuid"
)

func Login(reg *admin.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password := r.FormValue("email"), r.FormValue("password")
		var user models.AdminUser
		if err := reg.DB.Where("email = ?", email).First(&user).Error; err != nil {
			RenderLogin(reg, w, r, "Invalid credentials")
			return
		}
		if !user.CheckPassword(password) {
			RenderLogin(reg, w, r, "Invalid credentials")
			return
		}
		sessionID := uuid.New().String()
		reg.DB.Create(&models.Session{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Duration(reg.Config.SessionTTL) * time.Hour),
		})
		http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: sessionID, Path: "/admin", HttpOnly: true})
		reg.SetFlash(w, "Login successful! Welcome back.")
		http.Redirect(w, r, "/admin", 303)
	}
}

func Logout(reg *admin.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("admin_session")
		if cookie != nil {
			reg.DB.Delete(&models.Session{}, "id = ?", cookie.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "", Path: "/admin", Expires: time.Unix(0, 0), MaxAge: -1, HttpOnly: true})
		http.Redirect(w, r, "/admin/login", 303)
	}
}

func RenderLogin(reg *admin.Registry, w http.ResponseWriter, r *http.Request, errorMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFS(admin.TemplateFS, "templates/login.html"))
	styleContent, err := admin.TemplateFS.ReadFile("templates/style.css")
	if err != nil {
		styleContent = []byte("")
	}

	pd := view.PageData{
		SiteTitle: reg.Config.SiteTitle,
		Error:     errorMsg,
		CSS:       template.CSS(styleContent),
	}
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, "Template error", 500)
		return
	}
}
