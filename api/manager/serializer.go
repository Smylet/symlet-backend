package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/token"
)

type HostelManagerSerializer struct {
	UID           uuid.UUID     `json:"-"`
	hostelManager HostelManager `json:"-"`
}

func (m *HostelManagerSerializer) CreateTx(ctx *gin.Context, db *gorm.DB) error {

	authPayload := ctx.MustGet(users.AuthorizationPayloadKey).(*token.Payload)
	hostelManager := HostelManager{
		UserID: authPayload.UserID,
	}
	if err := db.Create(&hostelManager).Error; err != nil {
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
