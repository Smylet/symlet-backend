package db

import (
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		student.Student{},
		hostel.Hostel{},
		hostel.HostelStudent{},
	)
	if err != nil {
		panic(err)
	}
}
