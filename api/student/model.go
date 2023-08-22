package student

import "github.com/Smylet/symlet-backend/api/users"

// Student is a form of user model for our application
type Student struct {
	users.User `gorm:"embedded"`
	University string `gorm:"not null"`

}
