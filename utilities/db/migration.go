package db

import (
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		users.User{},
		users.Profile{},
		student.Student{},
		hostel.Hostel{},
		hostel.HostelStudent{},
	)
	if err != nil {
		panic(err)
	}
}
