package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/stdlib"
	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDBInstance is the Postgres-specific DbInstance variant.
type PostgresDBInstance struct {
	DBInstance
}

// Reset implementation for this type.
func (pgdb PostgresDBInstance) Reset() error {
	log.Info("Resetting database schema")
	if err := pgdb.GormDB().Exec("drop schema public cascade").Error; err != nil {
		return eris.Wrap(err, "error attempting to drop schema")
	}
	if err := pgdb.GormDB().Exec("create schema public").Error; err != nil {
		return eris.Wrap(err, "error attempting to create schema")
	}
	return nil
}

// NewPostgresDBInstance will construct a Postgres DbInstance.
func NewPostgresDBInstance(
	dsnURL string, slowThreshold time.Duration,
) (*PostgresDBInstance, error) {
	pgdb := PostgresDBInstance{
		DBInstance: DBInstance{dsn: dsnURL},
	}

	postgresConfig, err := DSNToConnConfig(dsnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	sqlDB := stdlib.OpenDB(postgresConfig)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour * 2)

	gormConfig := &gorm.Config{
		Logger: NewLoggerAdaptor(log.StandardLogger(), LoggerAdaptorConfig{
			SlowThreshold:             slowThreshold,
			IgnoreRecordNotFoundError: true,
		}),
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), gormConfig)
	if err != nil {
		log.Warnf("Failed to connect using provided host, retrying with localhost: %v", err)

		postgresConfig.Host = "localhost"
		sqlDB = stdlib.OpenDB(postgresConfig)
		gormDB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), gormConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pgdb.DB = gormDB
	return &pgdb, nil
}
