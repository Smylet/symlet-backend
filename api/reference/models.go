package reference

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Smylet/symlet-backend/utilities/common"
	"gorm.io/gorm"
)

type ReferenceModelInterface interface {
	common.ModelInterface
	Populate(db *gorm.DB) error
	GetTableName() string
}

type ReferenceHostelAmmenities struct {
	common.AbstractBaseReferenceModel
	Name        string
	Description string
}

func (h ReferenceHostelAmmenities) GetTableName() string {
	return "reference_hostel_ammenities"
}

func (h ReferenceHostelAmmenities) Populate(db *gorm.DB) error {
	// Populate the ammenities table with the data from the json file
	file, err := os.Open("../../resources/amenities.json")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	err = db.First(&ReferenceHostelAmmenities{}).Error
	if err == nil {
		log.Print("Ammenities table already populated")
		return nil
	}

	var amenities []*ReferenceHostelAmmenities

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&amenities); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	// Create records in the database for each amenity
	// BatchCreate
	for _, amenity := range amenities {
		result := db.Create(amenity)
		if result.Error != nil {
			fmt.Println("Error inserting data:", result.Error)
			return err
		}
		log.Printf("%v inserted\n", amenity.Name)
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
	file, err := os.Open("../../resources/universities.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = db.First(&ReferenceUniversity{}).Error
	if err == nil {
		log.Print("University table already populated")
		return nil
	}

	var universities []ReferenceUniversity
	err = json.NewDecoder(file).Decode(&universities)
	if err != nil {
		log.Fatal(err)
	}

	for _, uni := range universities {
		result := db.Create(&uni)
		if result.Error != nil {
			log.Println(result.Error)
		}
		log.Printf("%v inserted\n", uni.Name)

	}
	return nil
}

var ReferenceModelMap = map[string]ReferenceModelInterface{
	"amenities":  ReferenceHostelAmmenities{},
	"university": ReferenceUniversity{},
}
