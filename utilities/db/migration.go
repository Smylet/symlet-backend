package db

import (
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		users.User{},
		users.Profile{},
		users.VerificationEmail{},
		student.Student{},
		hostel.Hostel{},
		hostel.HostelStudent{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate")
	}
}
