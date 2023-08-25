package db

import (
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		student.Student{},
		hostel.Hostel{},
		hostel.HostelStudent{},
	)
}
