package models

// Permission defines what a role can do with a resource.
type Permission struct {
	ID           uint   `gorm:"primaryKey"`
	Role         string `gorm:"index"`
	ResourceName string `gorm:"index"`
	Action       string
}
