package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"gorm.io/gorm"
)

var (
    DB *gorm.DB
	initialized bool
)


func connectToPostgreSQL(config utils.Config, dbname string ) (*sql.DB, error) {
    // Define the PostgreSQL connection parameters

	postgresConfig := pgx.ConnConfig{
		Host:                 config.DBHost,
		Port:                 5432,
		User:                 config.DBUser,
		Password:             config.DBPass,
		Database:             dbname,
		PreferSimpleProtocol: true,
	}

	sqlDB := stdlib.OpenDB(postgresConfig)
	// Connection string
    


    return sqlDB, nil
}

func CreateTestDatabase(db *sql.DB) error {
    // Create the test database if not exist
    //-- Delete the database if it exists
    log.Println("Droping test database if it exist before creating")
    err := DropTestDatabase(db)
    if err != nil {
        return err
    }

    //DROP DATABASE IF EXISTS your_database_name;

    //-- Create the database if it doesn't exist
    log.Println("creating test database")
	res, err := db.Exec(
        "CREATE DATABASE smy_test",
	)
    log.Println(res)

	return err
}

func DropTestDatabase(db *sql.DB) error {
    dropTables :=fmt.Sprintln(`DO $$ 
    DECLARE 
        tabname text;
    BEGIN
        FOR tabname IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') 
        LOOP
            EXECUTE 'DROP TABLE IF EXISTS ' || tabname || ' CASCADE';
        END LOOP;
    END $$;`)
    _, err := db.Exec(dropTables)
    if err != nil {
        return err
    }
    _, err = db.Exec("DROP DATABASE IF EXISTS smy_test WITH (FORCE);")
    return err
}

func SetupTestDB() {
    if initialized {
        return
    }
    // Create a database
    log.Println(os.Getwd())
    config, err := utils.LoadConfig()
    if err != nil {
        panic("Failed to load config: " + err.Error())
    }

    sqlDB, err := connectToPostgreSQL(config, "")
    // Connect to the test database
    if err != nil {
        panic("Failed to open GORM database: " + err.Error())
    }

    // Create the test database
    if err := CreateTestDatabase(sqlDB); err != nil {
        panic("Failed to create test database: " + err.Error())
    }

    //sqlDB, err = connectToPostgreSQL("localhost", "postgres", "postgres", "smy_test")
    // Connect to the test database
    if err != nil {
        panic("Failed to open GORM database: " + err.Error())
    }
    if err != nil {
        panic("Failed to load config: " + err.Error())
    }
    config.DBName = "smy_test"
    DB = db.InitDB(config)
    if DB == nil {
        panic("Failed to initialize database")
    }
	// DB, err = gorm.Open(postgres.New(postgres.Config{
	// 	Conn: sqlDB,
	// }), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	// })
	// if err != nil {
	// 	log.Println("db err: (Init) ", err)
	// }

    // // Migrate the schema
    // DB.AutoMigrate(
    //     &reference.ReferenceHostelAmmenities{},
    //     &reference.ReferenceUniversity{},
    // )
    //populate Reference model
    for _, model := range reference.ReferenceModelMap {
        err := model.Populate(DB)
        if err != nil {
            panic("Failed to populate reference models: " + err.Error())
        }
    }

	initialized = true
}

func TeardownTestDB() {
    // Drop the test database
    if !initialized {
        return
    }
    config, err := utils.LoadConfig()
    if err != nil {
        panic("Failed to load config: " + err.Error())
    }

    sqlDB, err := connectToPostgreSQL(config, "")
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
    err := os.Setenv("ENV", "test")
    if err != nil {
        panic("Failed to set environment variable: " + err.Error())
    }

	SetupTestDB()
	defer TeardownTestDB()

	return m.Run()
}
