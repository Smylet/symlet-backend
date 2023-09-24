package db

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/jackc/pgx"
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
	dsnURL url.URL, slowThreshold time.Duration, poolMax int, reset bool,
) (*PostgresDBInstance, error) {
	pgdb := PostgresDBInstance{
		DBInstance: DBInstance{dsn: dsnURL.String()},
	}

	port, err := strconv.ParseUint(dsnURL.Port(), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	postgresConfig := pgx.ConnConfig{
		Host:                 dsnURL.Hostname(),
		Port:                 uint16(port),
		User:                 dsnURL.User.Username(),
		Database:             dsnURL.Path[1:],
		Password:             config.DBPass,
		PreferSimpleProtocol: true,
		LogLevel:             pgx.LogLevelWarn,
		RuntimeParams:        map[string]string{"search_path": "public"},
	}

	sqlDB := stdlib.OpenDB(postgresConfig)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour * 2)

	log.Infof("Using database %s", dsnURL.String())

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
		gormDB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), gormConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pgdb.DB = gormDB
	return &pgdb, nil
}
