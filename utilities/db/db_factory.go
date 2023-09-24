package db

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Smylet/symlet-backend/utilities/utils"
)

// MakeDBProvider will create a DbProvider of the correct type from the parameters.
// ParseDSN parses the provided DSN and returns a URL.
func ParseDSN(dsn string) (*url.URL, error) {
	return url.Parse(dsn)
}

// CreateDBProvider creates a DB provider based on the given URL.
func CreateDBProvider(dsnURL *url.URL, slowThreshold time.Duration, poolMax int, reset bool) (DBProvider, error) {
	return NewPostgresDBInstance(*dsnURL, slowThreshold)
}

// LoadConfig loads the application configuration.
func LoadConfig() (utils.Config, error) {
	return utils.LoadConfig()
}

// HandleResetAndMigration handles the reset and migration for the provided DB.
func HandleResetAndMigration(config utils.Config, db DBProvider) error {
	if config.Environment == "development" || config.Environment == "test" {
		if config.DatabaseReset {
			if err := db.Reset(); err != nil {
				return err
			}
		}

		return Migrate(db.GormDB(), config)
	}
	return nil
}

// MakeDBProvider orchestrates the above functions.
func MakeDBProvider(
	config utils.Config,
) (DBProvider, error) {
	dsnURL, err := ParseDSN(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}

	db, err := CreateDBProvider(dsnURL, time.Duration(config.DatabaseSlowThreshold.Seconds()), config.DatabasePoolMax, config.DatabaseReset)
	if err != nil {
		return nil, err
	}

	if err := HandleResetAndMigration(config, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
