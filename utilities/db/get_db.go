package db

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/Smylet/symlet-backend/utilities/utils"
)

// GetDB returns the reference to the database connection.
func GetDB() (DBProvider, error) {
	config := *utils.EnvConfig
	if config.Environment == "production" {
		config.DatabaseReset = false
	} else {
		config.DatabaseReset = true
	}

	// Log the connection info regardless of the environment
	str := fmt.Sprintf("Connecting to %s database", config.Environment)
	log.Info().Msg(str)

	database, err := MakeDBProvider(config)
	if err != nil {
		log.Error().Err(err).Msg("failed to get db")
		return nil, fmt.Errorf("error connecting to DB: %w", err)
	}

	return database, nil
}
