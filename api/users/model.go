// models.go
package users

import (
	"time"

	"github.com/Smylet/symlet-backend/api/core"
	"github.com/jinzhu/gorm"
)

type User struct {
	core.AbstractBaseModel
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `sql:"not null"`
	Profile      Profile
}

type Profile struct {
	ID       uint `gorm:"primary_key"`
	UserID   uint
	Username string
	Bio      string
	Image    string
}

var db *gorm.DB

func CreateUserTx() error {
	tx := db.Begin()

	// logic to create a user
	return tx.Commit().Error
}
