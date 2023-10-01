// models.go
package users

import (
	"time"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/google/uuid"
)

type User struct {
	common.AbstractBaseModel
	Username           string `gorm:"unique;not null"`
	Email              string `gorm:"unique;not null"`
	Password           string `sql:"not null"`
	Is_email_confirmed bool   `gorm:"default:false"`
}

type Profile struct {
	common.AbstractBaseModel
	UserID uint
	Bio    string
	Image  string
}

type VerificationEmail struct {
	ID         uint `gorm:"primary_key"`
	Email      string
	SecretCode string
	ExpiresAt  time.Time
	UserID     uint
}

type Session struct {
	ID           uuid.UUID
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiresAt    time.Time
	IsBlocked    bool
}
