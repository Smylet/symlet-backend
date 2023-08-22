package student

import (
	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/users"
)

// Student is a form of user model for our application
type Student struct {
	core.AbstractBaseModel
	users.User `gorm:"embedded"`
	University string `gorm:"not null"`
}
