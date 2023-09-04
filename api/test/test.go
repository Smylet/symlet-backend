package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	initialized bool
)

func connectToPostgreSQL(host, username, password string) (*sql.DB, error) {
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
	_, err := db.Exec(
		`CREATE DATABASE IF NOT EXIST smy_test ;`,
	)

	return err
}

func DropTestDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS smy_test;")
	return err
}

func SetupTestDB() {
	if initialized {
		return
	}
	// Create a database
	sqlDB, err := connectToPostgreSQL("localhost", "postgres", "postgres")
	// Connect to the test database
	if err != nil {
		panic("Failed to open GORM database: " + err.Error())
	}

	// Create the test database
	if err := CreateTestDatabase(sqlDB); err != nil {
		panic("Failed to create test database: " + err.Error())
	}
	config, err := utils.LoadConfig("../..")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	DB = db.InitDB(config)

	// populate Reference model
	for _, model := range reference.ReferenceModelMap {
		err := model.Populate(DB)
		if err != nil {
			panic("Failed to populate reference models: " + err.Error())
		}
	}

	// Insert a test user
	user := users.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	if err := DB.Create(&user).Error; err != nil {
		panic("Failed to insert test user: " + err.Error())
	}
	var university reference.ReferenceUniversity
	DB.Model(&reference.ReferenceUniversity{}).First(&university)

	student := student.Student{
		User:       user,
		University: university,
	}
	if err := DB.Create(&student).Error; err != nil {
		panic("Failed to insert test student: " + err.Error())
	}

	initialized = true
}

func TeardownTestDB() {
	// Drop the test database
	if !initialized {
		return
	}
	sqlDB, err := connectToPostgreSQL("localhost", "postgres", "postgres")
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
