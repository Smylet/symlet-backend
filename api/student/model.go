package student

import (
	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/users"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Student{})
// }

// Student is a form of user model for our application
type Student struct {
	core.AbstractBaseModel
	UserID	  uint   `gorm:"not null"`
	User  users.User `gorm:"foreignKey:UserID"`
	University string `gorm:"not null"`
}
