package admin

import (
	"time"
)

// AuditLog records every change made in the admin panel.
type AuditLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index"`
	UserEmail    string
	ResourceName string    `gorm:"index"`
	RecordID     string    `gorm:"index"`
	Action       string    // Create, Update, Delete
	Changes      string    // JSON or text description of changes
	CreatedAt    time.Time `gorm:"index"`
}

// RecordAction saves an audit entry to the database.
func (reg *Registry) RecordAction(user *AdminUser, resourceName, recordID, action, changes string) {
	log := AuditLog{
		UserID:       user.ID,
		UserEmail:    user.Email,
		ResourceName: resourceName,
		RecordID:     recordID,
		Action:       action,
		Changes:      changes,
		CreatedAt:    time.Now(),
	}
	reg.DB.Create(&log)
}
