package manager

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/users"
)

type HostelManagerSerializer struct {
	UID           uuid.UUID     `json:"-"`
	hostelManager HostelManager `json:"-"`
}

func (m *HostelManagerSerializer) Create(ctx *gin.Context, db *gorm.DB) error {
	var user users.User
	payload, err := users.GetAuthPayloadFromCtx(ctx)
	if err != nil {
		return  err
	}
	err = db.Model(&users.User{}).Where("id = ?", payload.UserID).First(&user).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}
	hostelManager := HostelManager{
		User: user,
	}
	if err = db.Create(&hostelManager).Error; err != nil {
		return err
	}
	m.hostelManager = hostelManager
	m.UID = hostelManager.UID
	return nil
}

func (m *HostelManagerSerializer) Response() map[string]interface{} {
	return map[string]interface{}{
		"uid": m.UID,
	}
}
