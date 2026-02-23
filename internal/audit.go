package internal

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"time"
)

func RecordAction(reg *admin.Registry, user *models.AdminUser, resName, recordID, action, changes string) {
	reg.DB.Create(&models.AuditLog{
		UserID: user.ID, UserEmail: user.Email, ResourceName: resName, 
		RecordID: recordID, Action: action, Changes: changes, CreatedAt: time.Now(),
	})
}
