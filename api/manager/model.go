package manager

import (
	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/users"
)


type HostelManager struct {
	core.AbstractBaseModel
	UserID uint `gorm:"not null"`
	User  users.User `gorm:"foreignKey:UserID"`
}