package models

import "time"

// AuditLog records every change made in the admin panel.
type AuditLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index"`
	UserEmail    string
	ResourceName string    `gorm:"index"`
	RecordID     string    `gorm:"index"`
	Action       string    
	Changes      string    
	CreatedAt    time.Time `gorm:"index"`
}
