// Package internal contains internal helpers for auth, CRUD and auditing.
package internal

import (
	"time"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
)

// RecordAction logs an audit record for a user action.
func RecordAction(reg *admin.Registry, user *models.AdminUser, resName, recordID, action, changes string) {
	reg.DB.Create(&models.AuditLog{
		UserID: user.ID, UserEmail: user.Email, ResourceName: resName,
		RecordID: recordID, Action: action, Changes: changes, CreatedAt: time.Now(),
	})
}
