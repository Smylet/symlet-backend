package student

import (
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Student{})
// }

// Student is a form of user model for our application
type Student struct {
	common.AbstractBaseModel
	users.User `gorm:"embedded"`
	University string `gorm:"not null"`
}
