package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Smylet/symlet-backend/utilities/utils"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database struct represents the database connection.
type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Init opens a database connection and saves the reference to the Database struct.
func InitDB(config utils.Config) *gorm.DB {
	// Adjust the connection string based on your PostgreSQL setup
	connectionString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=UTC",
		config.DBHost, config.DBUser, config.DBPass, config.DBName, config.DBPort, config.SSLMode,
	)

	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("db err: (Init) ", err)
	}

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
		log.Println("db err: (Init) ", err)
	}

	// Only migrate if their is a change in schema
	Migrate(db)

	return db
}

// GetDB returns the reference to the database connection.
func GetDB(config utils.Config) *gorm.DB {
	return InitDB(config)
}
