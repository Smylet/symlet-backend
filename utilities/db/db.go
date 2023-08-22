package db

import (
	"fmt"

	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database struct represents the database connection.
type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Init opens a database connection and saves the reference to the Database struct.
func InitDB() *gorm.DB {
	// Adjust the connection string based on your PostgreSQL setup
	connectionString := "host=localhost user=your_username dbname=your_dbname sslmode=disable password=your_password"
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}
	
	//Only migrate if their is a change in schema
	Migrate(db)

	db.DB().SetMaxIdleConns(10)
	//db.LogMode(true)
	DB = db
	return DB
}

// GetDB returns the reference to the database connection.
func GetDB() *gorm.DB {
	return DB
}


func GetModels() []core.ModelInterface {
	return []core.ModelInterface{
		student.Student{},
		hostel.Hostel{},
		hostel.HostelStudent{},
	}

}