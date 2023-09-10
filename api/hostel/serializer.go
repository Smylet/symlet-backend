package hostel

import (
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/token"
)

type HostelFeeSerializer struct {
	TotalAmount float64            `json:"total_amount" binding:"required"`
	PaymentPlan string             `json:"payment_plan"  binding:"required,oneof=monthly by_school_session annually"`
	Breakdown   cusjsonb `json:"breakdown" binding:"required"`
}

type AmmenitySerializer struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type HostelSerializer struct {
	ManagerID             uint   `json:"-" form:"-"`
	UniversityID          uint   `json:"university_id" form:"university_id" binding:"required"`
	Name                  string `json:"name" form:"name" binding:"required"`
	Address               string `json:"address" form:"address" binding:"required"`
	City                  string `json:"city" form:"city" binding:"required"`
	State                 string `json:"state" form:"state" binding:"required"`
	Country               string `json:"country" form:"country" binding:"required"`
	Description           string `json:"description" form:"discription" binding:"required"`
	NumberOfUnits         uint   `json:"number_of_units" binding:"required"`
	NumberOfOccupiedUnits uint   `json:"number_of_occupied_units" binding:"required"`
	NumberOfBedrooms      uint   `json:"number_of_bedrooms" binding:"required"`
	NumberOfBathrooms     uint   `json:"number_of_bathrooms" binding:"required"`
	Kitchen               string `json:"kitchen" binding:"required,oneof=shared none private"`

	FloorSpace uint                    `json:"floor_space"`
	HostelFee  HostelFeeSerializer     `json:"hostel_fee"`
	Amenities  []AmmenitySerializer    `json:"amenities"`
	Images     []*multipart.FileHeader `form:"images" binding:"max=10" swaggerignore:"true"`
	Hostel     Hostel                  `json:"-" swaggerignore:"true"`
}

func (h *HostelSerializer) AfterCreate() error {
	return nil
}

// CreateTx creates a new hostel
func (h *HostelSerializer) CreateTx(ctx *gin.Context, db *gorm.DB, session *session.Session) error {
	// Validate the fields
	if err := h.Validate(); err != nil {
		return err
	}
	// Get the manager id from the auth payload
	authPayload := ctx.MustGet(users.AuthorizationPayloadKey).(*token.Payload)
	var hostelManager manager.HostelManager

	err := db.Model(&manager.HostelManager{}).Where("user_id = ?", authPayload.UserID).First(&hostelManager).Error
	if err != nil {
		return fmt.Errorf("failed to find manager with user id %d: %v", authPayload.UserID, err)
	}
	h.ManagerID = hostelManager.ID

	// Process the uploaded images
	//filePaths, err := ProcessUploadedImages(h.Images, session)
	//if err != nil {
	//	return err
	//}
	// hostelImages := make([]HostelImage, len(filePaths))

	// for i, image := range filePaths {
	// 	hostelImages[i] = HostelImage{}
	// 	hostelImages[i].ImageURL = image
	// }
	// Create the hostel
	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		hostel := Hostel{
			ManagerID:             h.ManagerID,
			Name:                  h.Name,
			UniversityID:          h.UniversityID,
			Address:               h.Address,
			City:                  h.City,
			State:                 h.State,
			Country:               h.Country,
			Description:           h.Description,
			NumberOfUnits:         h.NumberOfUnits,
			NumberOfOccupiedUnits: h.NumberOfOccupiedUnits,
			NumberOfBedrooms:      h.NumberOfBedrooms,
			NumberOfBathrooms:     h.NumberOfBathrooms,
			Kitchen:               h.Kitchen,
			FloorSpace:            h.FloorSpace,
		}
		fmt.Println(hostel)

		if err := tx.Create(&hostel).Error; err != nil {
			return err
		}
		// Add the amenities
		amenitiesArr := make([]uint, len(h.Amenities))
		for _, ammenity := range h.Amenities {
			amenitiesArr = append(amenitiesArr, ammenity.ID)
		}
		amenities := []reference.ReferenceHostelAmenities{}
		if err := tx.Where("id IN ?", amenitiesArr).Find(&amenities).Error; err != nil {
			return err
		}
		fmt.Println(amenities)

		if err := tx.Model(&hostel).Association("Amenities").Append(amenities); err != nil {
			return err
		}
		// Add the hostel images
		// for _, image := range hostelImages {
		// 	image.HostelID = hostel.ID
		// }

		// if err := tx.Model((&HostelImage{})).Create(&hostelImages).Error; err != nil {
		// 	return err
		// }
		// Add the hostel fee
		// hostelFee := HostelFee{
		// 	HostelID:    hostel.ID,
		// 	TotalAmount: h.HostelFee.TotalAmount,
		// 	PaymentPlan: h.HostelFee.PaymentPlan,
		// 	Breakdown:   h.HostelFee.Breakdown,
		// }

		// if err := tx.Model((&HostelFee{})).Create(&hostelFee).Error; err != nil {
		// 	return err
		// }

		h.Hostel = hostel
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Validate validates the fields of the hostel
func (h *HostelSerializer) Validate() error {
	//Validate the fields
	var errorMessage string
	validate := validator.New()
	if err := validate.Struct(h); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Return validation errors to the client
		for _, err := range validationErrors {
			errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
		}
		return fmt.Errorf(errorMessage)
	}
	return nil
}

// Response returns the response of the hostel
func (h HostelSerializer) Response() map[string]interface{} {
	amenities := make([]map[string]interface{}, len(h.Amenities))
	for i, ammenity := range h.Amenities {
		amenities[i] = map[string]interface{}{
			"id":   ammenity.ID,
			"name": ammenity.Name,
		}
	}
	imageURLs := make([]string, len(h.Hostel.HostelImages))
	for i, image := range h.Hostel.HostelImages {
		imageURLs[i] = image.ImageURL
	}
	breakDown := make(map[string]float64, len(h.Hostel.HostelFee.Breakdown))
	for key, value := range h.Hostel.HostelFee.Breakdown {
		breakDown[key] = value
	}

	fmt.Println(h.Hostel.HostelFee)

	hostelFee := map[string]interface{}{
		"total_amount": h.Hostel.HostelFee.TotalAmount,
		"payment_plan": h.Hostel.HostelFee.PaymentPlan,
		"breakdown":    breakDown,
	}

	return map[string]interface{}{
		"uid":                      h.Hostel.UID,
		"manager_uid":              h.Hostel.Manager.UID,
		"name":                     h.Hostel.Name,
		"university":               h.Hostel.University.Name,
		"address":                  h.Hostel.Address,
		"city":                     h.Hostel.City,
		"state":                    h.Hostel.State,
		"country":                  h.Hostel.Country,
		"description":              h.Hostel.Description,
		"number_of_units":          h.Hostel.NumberOfUnits,
		"number_of_occupied_units": h.Hostel.NumberOfOccupiedUnits,
		"number_of_bedrooms":       h.Hostel.NumberOfBedrooms,
		"number_of_bathrooms":      h.Hostel.NumberOfBathrooms,
		"kitchen":                  h.Hostel.Kitchen,
		"floor_space":              h.Hostel.FloorSpace,
		"hostel_fee":               hostelFee,
		"amenities":                amenities,
		"hostel_images":            imageURLs,
	}
}

func (h *HostelSerializer) GetHostelTx(db *gorm.DB, hostelUID uuid.UUID) error {
	err := db.Model(&Hostel{}).Preload(clause.Associations).Where("uid = ?", hostelUID).First(&h.Hostel).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *HostelSerializer) ListHostelsTx(db *gorm.DB, queryParams HostelQueryParams) ([]Hostel, error) {
	hostels, err := h.FilterHostels(db, queryParams)
	if err != nil {
		return nil, err
	}

	return hostels, nil
}

func (h HostelSerializer) ResponseMany(hostels []Hostel) []map[string]interface{} {
	hostelsResponse := make([]map[string]interface{}, len(hostels))
	for i, hostel := range hostels {
		hostelsResponse[i] = HostelSerializer{Hostel: hostel}.Response()
	}
	return hostelsResponse
}

func (h *HostelSerializer) UpdateHostelTx(ctx *gin.Context, db *gorm.DB, session *session.Session, hostelUID uuid.UUID) error {
	var hostel Hostel

	if err := h.Validate(); err != nil {
		return err
	}
	authPayload := ctx.MustGet(users.AuthorizationPayloadKey).(*token.Payload)

	err := db.Model(&manager.HostelManager{}).Where("user_id = ?", authPayload.UserID).First(&h.ManagerID).Error
	if err != nil {
		return fmt.Errorf("failed to find manager with user id %d: %v", authPayload.UserID, err)
	}
	if h.ManagerID != hostel.ManagerID {
		return fmt.Errorf("user is not authorized to update this hostel")
	}

	updatedFields := h.getUpdatedFields()

	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		err := tx.Model(&Hostel{}).Preload("Amenities").Where("uid = ?", hostelUID).First(&hostel).Error
		if err != nil {
			return err
		}

		if err := tx.Model(&hostel).Updates(updatedFields).Error; err != nil {
			return err
		}
		if h.Amenities != nil {
			amenitiesArr, err := h.updateAmenities(&hostel, tx)
			if err != nil {
				return err
			}

			amenities := []reference.ReferenceHostelAmenities{}
			if err := tx.Where("id IN ?", amenitiesArr).Find(&amenities).Error; err != nil {
				return err
			}

			if err := tx.Model(&hostel).Association("Amenities").Append(amenities); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *HostelSerializer) updateAmenities(hostel *Hostel, tx *gorm.DB) ([]uint, error) {
	amenitiesArr := make([]uint, len(h.Amenities))

	// Remove amenities that are not in the new list
	for _, amenity := range hostel.Amenities {
		found := false
		for _, newAmenity := range h.Amenities {
			if amenity.ID == newAmenity.ID {
				found = true
				break
			}
		}
		if !found {
			if err := tx.Model(&hostel).Association("Amenities").Delete(&amenity); err != nil {
				return nil, err
			}
		}
	}

	// Add amenities that are not in the old list
	for _, newAmenity := range h.Amenities {
		found := false
		for _, amenity := range hostel.Amenities {
			if amenity.ID == newAmenity.ID {
				found = true
				break
			}
		}
		if !found {
			amenitiesArr = append(amenitiesArr, newAmenity.ID)
		}
	}
	return amenitiesArr, nil
}

// Function to filter hostels based on query parameters
func (h *HostelSerializer) FilterHostels(db *gorm.DB, queryParams HostelQueryParams) ([]Hostel, error) {
	// Initialize the query with the Hostel model
	var hostels []Hostel

	query := db.Model(&Hostel{})

	// Apply filters based on the queryParams fields
	// University ID Filter
	if queryParams.UniversityID != 0.0 {
		query = query.Where("university_id = ?", queryParams.UniversityID)
	}

	// Hostel Fee Filters
	if queryParams.HostelFeeTotalMin != 0.0 {
		query = query.Where("hostel_fee_total >= ?", queryParams.HostelFeeTotalMin)
	}
	if queryParams.HostelFeeTotalMax != 0.0 {
		query = query.Where("hostel_fee_total <= ?", queryParams.HostelFeeTotalMax)
	}
	if queryParams.HostelFeePlan != "" {
		query = query.Where("hostel_fee_plan = ?", queryParams.HostelFeePlan)
	}

	// Amenities Filter
	if queryParams.HasAmenities {
		if queryParams.HasAmenities {
			query = query.Where("amenities IS NOT NULL")
		} else {
			query = query.Where("amenities IS NULL")
		}
	}

	if queryParams.Name != "" {
		query = query.Where("name LIKE ?", "%"+queryParams.Name+"%")
	}

	// Address starts with filter
	if queryParams.Address != "" {
		query = query.Where("address LIKE ?", queryParams.Address+"%")
	}

	// Description ends with filter
	if queryParams.Description != "" {
		query = query.Where("description LIKE ?", "%"+queryParams.Description)
	}

	// City filter
	if queryParams.City != "" {
		query = query.Where("city = ?", queryParams.City)
	}

	// Security Rating Filters
	if queryParams.SecurityRatingMin != nil {
		query = query.Where("security_rating >= ?", *queryParams.SecurityRatingMin)
	}
	if queryParams.SecurityRatingMax != nil {
		query = query.Where("security_rating <= ?", *queryParams.SecurityRatingMax)
	}

	// Location Rating Filters
	if queryParams.LocationRatingMin != nil {
		query = query.Where("location_rating >= ?", *queryParams.LocationRatingMin)
	}
	if queryParams.LocationRatingMax != nil {
		query = query.Where("location_rating <= ?", *queryParams.LocationRatingMax)
	}

	// General Rating Filters
	if queryParams.GeneralRatingMin != nil {
		query = query.Where("general_rating >= ?", *queryParams.GeneralRatingMin)
	}
	if queryParams.GeneralRatingMax != nil {
		query = query.Where("general_rating <= ?", *queryParams.GeneralRatingMax)
	}

	// Continue with similar logic for other rating fields...

	// You can also apply other filters such as name, address, etc., as needed
	err := query.Find(&hostels).Error
	if err != nil {
		return nil, err
	}

	return hostels, nil
}

// // Retrieves the fields sent in the request body of the update request
func (d *HostelSerializer) getUpdatedFields() map[string]interface{} {
	data := map[string]interface{}{
		"name":                d.Name,
		"university_id":       d.UniversityID,
		"address":             d.Address,
		"city":                d.City,
		"state":               d.State,
		"country":             d.Country,
		"description":         d.Description,
		"number_of_units":     d.NumberOfUnits,
		"number_of_bedrooms":  d.NumberOfBedrooms,
		"number_of_bathrooms": d.NumberOfBathrooms,
		"kitchen":             d.Kitchen,
		"floor_space":         d.FloorSpace,
	}
	for key, value := range data {
		val := reflect.ValueOf(value)
		if !val.IsValid() || reflect.DeepEqual(value, reflect.Zero(val.Type()).Interface()) {
			delete(data, key)
		}
	}
	return data

}
