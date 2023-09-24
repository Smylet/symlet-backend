package db

import (
	"fmt"

	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/rs/zerolog/log"
)

// GetDB returns the reference to the database connection.
func GetDB(config utils.Config) (DBProvider, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msg("failed to get db")
		}
	}()

	if config.Environment == "production" {
		config.DatabaseReset = false

		// connected to config.DatabaseURI
		str := fmt.Sprintf("Connecting to %s database", config.Environment)
		log.Info().Msg(str)

	} else {
		config.DatabaseReset = true

		// connected to config.DatabaseURI
		str := fmt.Sprintf("Connecting to %s database", config.Environment)
		log.Info().Msg(str)
	}

	database, err := MakeDBProvider(
		config.DatabaseURI,
		config.DatabaseSlowThreshold,
		config.DatabasePoolMax,
		config.DatabaseReset,
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %w", err)
	}
	// cache a global reference to the gorm.DB
	DB = database.GormDB()
	return database, nil
}
