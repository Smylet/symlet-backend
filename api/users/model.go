// models.go
package users

import (
	"time"

	"github.com/google/uuid"

	"github.com/Smylet/symlet-backend/utilities/common"
)

type User struct {
	common.AbstractBaseModel
	UserName         string `gorm:"unique;not null"`
	Email            string `gorm:"unique;not null"`
	Password         string `sql:"not null"`
	IsEmailConfirmed bool   `gorm:"default:false"`

	// Polymorphic relationship with Student, Vendor and Hostel Owner
	RoleID   uint
	RoleType string

	ProfileID uint
	Profile   Profile `gorm:"foreignKey:UserID;references:id"`

	ResetPasswordToken     string
	ResetPasswordExpiresAt time.Time
	EmailReminderCount     uint

	IsTwoFAEnabled bool `gorm:"default:false"`
	TwoFASecret    uint

	Preferences  string `gorm:"type:json;default:'{}'"` // serialized JSON for preferences
	PastSearches string `gorm:"type:json;default:'{}'"`
}

type Profile struct {
	common.AbstractBaseModel
	UserID    uint
	Bio       string
	Image     string
	FirstName string
	LastName  string
}

type VerificationEmail struct {
	ID         uint `gorm:"primary_key"`
	Email      string
	SecretCode uint
	ExpiresAt  time.Time
	UserID     uint
}

type Session struct {
	ID           uuid.UUID
	UserID       uint
	UserName     string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiresAt    time.Time
	IsBlocked    bool
}
