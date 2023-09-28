package reference

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)


type AmenitySerializer struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}


func (serializer *AmenitySerializer) List(db *gorm.DB) ([]ReferenceHostelAmenities, error) {
	var amenities []ReferenceHostelAmenities
	err := db.Find(&amenities).Error
	if err != nil {
		return nil, err
	}
	return amenities, nil
}


func (serializer *AmenitySerializer) ResponseMany(amenities []ReferenceHostelAmenities) []AmenitySerializer {
	var response []AmenitySerializer
	for _, amenity := range amenities {
		response = append(response, AmenitySerializer{
			ID: amenity.ID,
			Name: amenity.Name,
		})
	}
	return response
}


type UniversitySerializer struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	City string `json:"city"`
	Country string `json:"country"`
	State string `json:"state"`
}



func (serializer *UniversitySerializer) List(db *gorm.DB, queryParams UniversityQueryParams) ([]ReferenceUniversity, error) {
	var universities []ReferenceUniversity
	universities, err := serializer.Filter(db, queryParams)
	if err != nil {
		return nil, err
	}
	return universities, nil
}


func (serializer *UniversitySerializer) ResponseMany(universities []ReferenceUniversity) []UniversitySerializer {
	var response []UniversitySerializer
	for _, university := range universities {
		response = append(response, UniversitySerializer{
			ID: university.ID,
			Name: university.Name,
			Code: university.Code,
			City: university.City,
			Country: university.Country,
			State: university.State,
		})
	}
	return response
}

func (serializer *UniversitySerializer)Filter(db *gorm.DB, queryParams UniversityQueryParams) ([]ReferenceUniversity, error) {
	var universities []ReferenceUniversity
	fmt.Print(queryParams)
	query := db.Model(&ReferenceUniversity{})

	// Case insensitive filter
	switch {
	case queryParams.Name != "":
		query = query.Where("lower(name) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(queryParams.Name)))
	case queryParams.Code != "":
		query = query.Where("code LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(queryParams.Code)))
	case queryParams.City != "":
		query = query.Where("lower(city) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(queryParams.City)))
	case queryParams.State != "":
		query = query.Where("lower(state) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(queryParams.State)))
	case queryParams.Country != "":
		query = query.Where("lower(country) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(queryParams.Country)))

	}

	err := query.Find(&universities).Error
	if err != nil {
		return nil, err
	}
	return universities, nil
}