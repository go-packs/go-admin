package internal

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"net/http"
	"time"
)

func IsAllowed(reg *admin.Registry, role, resource, action string) bool {
	if role == "admin" { return true }
	var count int64
	reg.DB.Model(&models.Permission{}).Where("role = ? AND resource_name = ? AND action = ?", role, resource, action).Count(&count)
	return count > 0
}

func GetUserFromRequest(reg *admin.Registry, r *http.Request) (*models.AdminUser, string) {
	cookie, err := r.Cookie("admin_session")
	if err != nil { return nil, "guest" }
	var sess models.Session
	if err := reg.DB.Where("id = ? AND expires_at > ?", cookie.Value, time.Now()).First(&sess).Error; err != nil { return nil, "guest" }
	var user models.AdminUser
	if err := reg.DB.First(&user, sess.UserID).Error; err != nil { return nil, "guest" }
	return &user, user.Role
}
