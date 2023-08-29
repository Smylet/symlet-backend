package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/utilities/utils"
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
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}

	// Only migrate if their is a change in schema
	Migrate(db)

	sqlDB, _ := db.DB()

	DB = db

	sqlDB.SetMaxIdleConns(10)

	//defer sqlDB.Close()
	// db.LogMode(true)
	return DB
}

// GetDB returns the reference to the database connection.
func GetDB(config utils.Config) *gorm.DB {
	return InitDB(config)
}

// var implementedModelInterface []core.ModelInterface

// //All package having a model should call this function in an init
// // function in the model.go file to register their model
// func RegisterModel(model ...core.ModelInterface) {
// 	fmt.Println(implementedModelInterface, "WFNORR")
// 	implementedModelInterface = append(implementedModelInterface, model...)
// }

// func GetModels() []core.ModelInterface{
// 	return implementedModelInterface
// }
