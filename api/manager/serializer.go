package manager

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	err = db.Model(&users.User{}).Preload(clause.Associations).Where("id = ?", payload.UserID).First(&user).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}

	if user.RoleType != ""{
		return fmt.Errorf("user is as already associated with a %v role", user.RoleType)
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
		"user": map[string]interface{}{
			"uid":        m.hostelManager.User.UID,
			"username":   m.hostelManager.User.Username,
			"email":      m.hostelManager.User.Email,
			"created_at": m.hostelManager.User.CreatedAt,
			"updated_at": m.hostelManager.User.UpdatedAt,
		},
		"profile": map[string]interface{}{
			"uid":        m.hostelManager.User.Profile.UID,
			"first_name": m.hostelManager.User.Profile.FirstName,
			"last_name":  m.hostelManager.User.Profile.LastName,
			"bio":        m.hostelManager.User.Profile.Bio,
			"image":      m.hostelManager.User.Profile.Image,
			"created_at": m.hostelManager.User.Profile.CreatedAt,
			"updated_at": m.hostelManager.User.Profile.UpdatedAt,
		},
	}
	}

