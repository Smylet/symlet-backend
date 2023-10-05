package manager

import (
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
)

type HostelManager struct {
	common.AbstractBaseModel
	User   users.User `gorm:"polymorphic:User"`
}
