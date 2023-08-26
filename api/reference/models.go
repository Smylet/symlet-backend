package reference

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Smylet/symlet-backend/api/core"
	"gorm.io/gorm"
)

type ReferenceModelInterface interface {
	core.ModelInterface
	isReferenceModel() bool
	Populate()
}

type ReferenceHostelAmmenities struct {
	core.AbstractBaseReferenceModel
	Name string
	Description string
}

func (h *ReferenceHostelAmmenities)Populate(db *gorm.DB) error{

	//Populate the ammenities table with the data from the json file
	file, err := os.Open("ammenities.json")
	if err != nil{
		fmt.Printf("Error opening file: %v", err)
		return err 
	}
	defer file.Close()


	var amenities []ReferenceHostelAmmenities
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&amenities); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}
	
	// Create records in the database for each amenity
	// BatchCreate

	result := db.Create(&amenities)
	if result.Error != nil {
		fmt.Println("Error inserting data:", result.Error)
		return err
	}
	return nil


}

type ReferenceUniversity struct {
	core.AbstractBaseReferenceModel
	Name string
	State string
	City string
	Country string
}	

func (u *ReferenceUniversity)Populate(db *gorm.DB) error{
	file, err := os.Open("universities.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
	}
	return nil
}



func Populate(model ReferenceModelInterface){
	model.Populate()	
}