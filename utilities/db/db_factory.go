package db

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/rotisserie/eris"
)

// MakeDBProvider will create a DbProvider of the correct type from the parameters.
func MakeDBProvider(
	dsn string, slowThreshold time.Duration, poolMax int, reset bool,
) (db DBProvider, err error) {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}
	switch dsnURL.Scheme {
	case "sqlite":
		db, err = NewSqliteDBInstance(
			*dsnURL,
			slowThreshold,
			poolMax,
			reset,
		)
		if err != nil {
			return nil, eris.Wrap(err, "error creating sqlite provider")
		}
	case "postgres", "postgresql":
		db, err = NewPostgresDBInstance(
			*dsnURL,
			slowThreshold,
			poolMax,
			reset,
		)
		if err != nil {
			return nil, eris.Wrap(err, "error creating postgres provider")
		}
	default:
		{
			return nil, eris.New("unsupported database type")
		}
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if config.Environment == "development" || config.Environment == "test" {
		if reset {
			if err := db.Reset(); err != nil {
				db.Close()
				return nil, err
			}
		}

		if err := Migrate(db.GormDB()); err != nil {
			db.Close()
			return nil, err
		}
	}

	// if err := checkAndMigrate(migrate, db); err != nil {
	// 	db.Close()
	// 	return nil, err
	// }

	return db, nil
}
