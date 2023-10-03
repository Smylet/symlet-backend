package reference

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

var logger = common.NewLogger()

type ReferenceModelInterface interface {
	common.ModelInterface
	Populate(db *gorm.DB) error
	GetTableName() string
}

type ReferenceHostelAmenities struct {
	common.AbstractBaseReferenceModel
	Name        string
	Description string
}

func (h ReferenceHostelAmenities) GetTableName() string {
	return "reference_hostel_ammenities"
}

func (h ReferenceHostelAmenities) Populate(db *gorm.DB) error {
	// Populate the ammenities table with the data from the json file
	config, err := utils.LoadConfig()
	if err != nil {
		logger.Error("Error loading config: ", err)
		return err
	}

	file, err := os.Open(filepath.Clean(config.BasePath) + "/resources/amenities.json")
	if err != nil {
		logger.Error("Error opening file: ", err)
		return err
	}
	defer file.Close()

	err = db.First(&ReferenceHostelAmenities{}).Error
	if err == nil {
		logger.Print(zerolog.InfoLevel, "Ammenities table already populated")
		return nil
	}

	var amenities []*ReferenceHostelAmenities

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&amenities); err != nil {
		logger.Error("Error decoding JSON", err) // Println("Error decoding JSON:", err)
		return err
	}

	// Create records in the database for each amenity
	// BatchCreate
	for _, amenity := range amenities {
		result := db.Create(amenity)
		if result.Error != nil {
			logger.Error("Error inserting data:", result.Error)
			return err
		}
		logger.Printf(context.Background(), "%v inserted\n", amenity.Name)
	}
	return nil
}

type ReferenceUniversity struct {
	common.AbstractBaseReferenceModel
	Name    string
	Slug    string
	State   string
	City    string
	Country string
	Code    string
}

func (u ReferenceUniversity) GetTableName() string {
	return "reference_university"
}

func (u ReferenceUniversity) Populate(db *gorm.DB) error {
	config, err := utils.LoadConfig()
	if err != nil {
		logger.Error(err)
		return err
	}

	file, err := os.Open(filepath.Clean(config.BasePath) + "/resources/universities.json")
	if err != nil {
		logger.Error(err)
	}
	defer file.Close()

	err = db.First(&ReferenceUniversity{}).Error
	if err == nil {
		logger.Print(zerolog.InfoLevel, "University table already populated")
		return nil
	}

	var universities []ReferenceUniversity
	err = json.NewDecoder(file).Decode(&universities)
	if err != nil {
		logger.Error(err)
	}

	for _, uni := range universities {
		result := db.Create(&uni)
		if result.Error != nil {
			logger.Error(result.Error)
		}
		logger.Printf(context.Background(), "%v inserted\n", uni.Name)

	}
	return nil
}

var ReferenceModelMap = map[string]ReferenceModelInterface{
	"amenities":  ReferenceHostelAmenities{},
	"university": ReferenceUniversity{},
}
