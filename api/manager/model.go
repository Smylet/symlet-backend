package manager

import (
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
)

type HostelManager struct {
	common.AbstractBaseModel
	UserID uint       `gorm:"not null"`
	User   users.User `gorm:"foreignKey:UserID"`
}
