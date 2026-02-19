package admin

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// AdminUser represents a person who can log in to the admin panel.
type AdminUser struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string
}

// Session stores active login sessions.
type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	ExpiresAt time.Time `gorm:"index"`
}

// SetPassword hashes and sets the user's password.
func (u *AdminUser) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies the password against the hash.
func (u *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Permission defines what a role can do with a resource.
type Permission struct {
	ID           uint   `gorm:"primaryKey"`
	Role         string `gorm:"index"`
	ResourceName string `gorm:"index"`
	Action       string
}

// IsAllowed checks if a role has permission for an action on a resource.
func (reg *Registry) IsAllowed(role string, resource string, action string) bool {
	if role == "admin" {
		return true
	}
	var count int64
	reg.DB.Model(&Permission{}).
		Where("role = ? AND resource_name = ? AND action = ?", role, resource, action).
		Count(&count)
	return count > 0
}

// GetUserFromRequest extracts the current session's user.
func (reg *Registry) GetUserFromRequest(r *http.Request) (*AdminUser, string) {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return nil, "guest"
	}

	var sess Session
	if err := reg.DB.Where("id = ? AND expires_at > ?", cookie.Value, time.Now()).First(&sess).Error; err != nil {
		return nil, "guest"
	}

	var user AdminUser
	if err := reg.DB.First(&user, sess.UserID).Error; err != nil {
		return nil, "guest"
	}

	return &user, user.Role
}
