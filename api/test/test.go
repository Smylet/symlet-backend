package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/Smylet/symlet-backend/api/reference"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	initialized bool
)


func connectToPostgreSQL(host, username, password, dbname string) (*sql.DB, error) {
    // Define the PostgreSQL connection parameters
    port := 5432

	// Connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, username, password,
	)

    // Open a connection to PostgreSQL
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func CreateTestDatabase(db *sql.DB) error {
    // Create the test database if not exist
    //-- Delete the database if it exists
    err := DropTestDatabase(db)
    if err != nil {
        return err
    }

    //DROP DATABASE IF EXISTS your_database_name;

    //-- Create the database if it doesn't exist

	_, err = db.Exec(
        "CREATE DATABASE smy_test",
	)

	return err
}

func DropTestDatabase(db *sql.DB) error {
    _, err := db.Exec("DROP DATABASE IF EXISTS smy_test WITH (FORCE);")
    return err
}

func SetupTestDB() {
    if initialized {
        return
    }
    // Create a database
    sqlDB, err := connectToPostgreSQL("localhost", "postgres", "postgres", "")
    // Connect to the test database
    if err != nil {
        panic("Failed to open GORM database: " + err.Error())
    }

    // Create the test database
    if err := CreateTestDatabase(sqlDB); err != nil {
        panic("Failed to create test database: " + err.Error())
    }
	logLevel := logger.Silent

    sqlDB, err = connectToPostgreSQL("localhost", "postgres", "postgres", "smy_test")
    // Connect to the test database
    if err != nil {
        panic("Failed to open GORM database: " + err.Error())
    }

	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		log.Println("db err: (Init) ", err)
	}

    // Migrate the schema
    DB.AutoMigrate(
        &reference.ReferenceHostelAmmenities{},
        &reference.ReferenceUniversity{},
    )
    //populate Reference model
    for _, model := range reference.ReferenceModelMap {
        err := model.Populate(DB)
        if err != nil {
            panic("Failed to populate reference models: " + err.Error())
        }
    }
    initialized = true

	initialized = true
}

func TeardownTestDB() {
    // Drop the test database
    if !initialized {
        return
    }
    sqlDB, err := connectToPostgreSQL("localhost", "postgres", "postgres", "")
    if err != nil {
        panic("Failed to open GORM database: " + err.Error())
    }
    if err := DropTestDatabase(sqlDB); err != nil {
        panic("Failed to drop test database: " + err.Error())
    }

    initialized = false
}

// RunTests is a helper function to run tests and handle setup/teardown.
func RunTests(m *testing.M) int {
	SetupTestDB()
	defer TeardownTestDB()

	return m.Run()
}
