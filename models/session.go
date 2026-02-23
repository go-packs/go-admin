package models

import "time"

// Session stores active login sessions.
type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	ExpiresAt time.Time `gorm:"index"`
}
