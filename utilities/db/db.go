package db

import (
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Smylet/symlet-backend/utilities/utils"
)

// Database struct represents the database connection.
type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Init opens a database connection and saves the reference to the Database struct.
func InitDB(config utils.Config) *gorm.DB {
	port, err := strconv.ParseUint(config.DBPort, 10, 16)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse db port")
	}

	postgresConfig := pgx.ConnConfig{
		Host:                 config.DBHost,
		Port:                 uint16(port),
		User:                 config.DBUser,
		Password:             config.DBPass,
		Database:             config.DBName,
		PreferSimpleProtocol: true,
	}

	sqlDB := stdlib.OpenDB(postgresConfig)

	maxIdleConns := 10               // Suitable for medium-sized applications.
	maxOpenConns := 20               // Depending on your expected load and DB capacity.
	connMaxLifetime := time.Hour * 2 // Recycle connections every 2 hours.
	logLevel := logger.Silent

	sqlDB.SetMaxIdleConns(maxIdleConns)       // max number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(maxOpenConns)       // max number of open connections in the database
	sqlDB.SetConnMaxLifetime(connMaxLifetime) // max amount of time a connection may be reused

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Only migrate if their is a change in schema - development or test mode
	// if config.Environment == "development" || config.Environment == "test" {
	// 	Migrate(config)
	// }

	return db
}

// GetDB returns the reference to the database connection.
func GetDB(config utils.Config) *gorm.DB {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msg("failed to get db")
		}
	}()
	return InitDB(config)
}
