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
	"github.com/Smylet/symlet-backend/utilities/utils"
)

type HostelFeeSerializer struct {
	TotalAmount float64 `json:"total_amount"  form:"total_amount"`
	PaymentPlan string  `json:"payment_plan" form:"payment_plan"` //binding:"oneof=monthly by_school_session annually"`
	Breakdown   Map     `json:"breakdown" form:"breakdown"`
}

type AmmenitySerializer struct {
	ID   uint   `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

type HostelSerializer struct {

	ManagerID             uint    `json:"-" form:"-"`
	UniversityID          uint    `json:"university_id" form:"university_id" custom_binding:"requiredForCreate"`
	Name                  *string `json:"name" form:"name" custom_binding:"requiredForCreate"`
	Address               *string `json:"address" form:"address" custom_binding:"requiredForCreate"`
	City                  *string `json:"city" form:"city" custom_binding:"requiredForCreate"`
	State                 *string `json:"state" form:"state" custom_binding:"requiredForCreate"`
	Country               *string `json:"country" form:"country" custom_binding:"requiredForCreate"`
	Description           *string `json:"description" form:"description" custom_binding:"requiredForCreate"`
	NumberOfUnits         *uint   `json:"number_of_units" form:"number_of_units" custom_binding:"requiredForCreate"`
	NumberOfOccupiedUnits *uint   `json:"number_of_occupied_units" form:"number_of_occupied_units" custom_binding:"requiredForCreate"`
	NumberOfBedrooms      *uint   `json:"number_of_bedrooms" form:"number_of_bedrooms" custom_binding:"requiredForCreate"`
	NumberOfBathrooms     *uint   `json:"number_of_bathrooms" form:"number_of_bathrooms" custom_binding:"requiredForCreate"`
	Kitchen               *string `json:"kitchen" form:"kitchen" custom_binding:"requiredForCreate" binding:"oneof=shared none private"`

	FloorSpace *uint                   `json:"floor_space" form:"floor_space" custom_binding:"requiredForCreate"`
	HostelFee  HostelFeeSerializer     `json:"hostel_fee" form:"hostel_fee"` //binding:"required"`
	Amenities  []AmmenitySerializer    `json:"amenities" form:"amenities"`   //binding:"required"`
	Images     []*multipart.FileHeader `form:"images" binding:"max=10" swaggerignore:"true" validate:"ValidateImageExtension"`
	Hostel     Hostel                  `json:"-" swaggerignore:"true"`
}

func (h *HostelSerializer) AfterCreate() error {
	return nil
}

func getManager(ctx *gin.Context, db *gorm.DB) (*manager.HostelManager, error) {
	// Get the manager id from the auth payload
	payload, err := users.GetAuthPayloadFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var hostelManager manager.HostelManager
	user := users.User{}

	err = db.Model(users.User{}).Where("id = ? AND role_type = ?", payload.UserID, "hostel_managers").First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user with id %d and role_type hostel_managers: %v", payload.UserID, err)
	}

	err = db.Model(&manager.HostelManager{}).Where("id = ?", user.RoleID).First(&hostelManager).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find manager with user id %d: %w", payload.UserID, err)
	}
	return &hostelManager, nil

}

// // CreateTx creates a new hostel
// func (h *HostelSerializer) CreateTx(ctx *gin.Context, db *gorm.DB, session *session.Session) error {
// 	logger := common.NewLogger()
// 	var wg sync.WaitGroup

// 	//Process the uploaded images

// 	// Create a channel to send errors from Goroutine to main thread
// 	errorChan := make(chan error, 1) // Buffer size 1 for single error

// 	// Create a channel to send file paths from Goroutine to main thread
// 	filePathsChan := make(chan string, len(h.Images))
// 	wg.Add(1)
// 	// Launch a Goroutine to process images concurrently
// 	go func() {
// 		defer wg.Done()
// 		logger.Info("Processing images in goroutine")

// 		// Process the uploaded images
// 		filePaths, err := ProcessUploadedImages(h.Images, session)
// 		if err != nil {
// 			logger.Info("Image processing failed")
// 			// Send the error to the main thread via the error channel
// 			errorChan <- err
// 			return
// 		}
// 		logger.Info(filePaths)

// 		// Send the file paths to the main thread via the channel
// 		for _, path := range filePaths {
// 			filePathsChan <- path
// 		}

// 		// Close the channels to signal that processing is complete
// 		close(errorChan)
// 		close(filePathsChan)
// 		logger.Info("Image processing complete")
// 		}()

// 	wg.Wait()
// 	hostelImagesSlice := make([]string, 0)
// 	logger.Info("Getting images from channel")
// 	for path := range filePathsChan {
// 		logger.Printf(ctx, "------> %v", path)
// 		hostelImagesSlice = append(hostelImagesSlice, path)
// 	}

// 	// Validate the fields
// 	// if err := h.Validate(); err != nil {
// 	// 	return err
// 	// }
// 	// Get the manager id from the auth payload
// 	authPayload := ctx.MustGet(users.AuthorizationPayloadKey).(*token.Payload)
// 	var hostelManager manager.HostelManager

// 	err := db.Model(&manager.HostelManager{}).Where("user_id = ?", authPayload.UserID).First(&hostelManager).Error
// 	if err != nil {
// 		return fmt.Errorf("failed to find manager with user id %d: %v", authPayload.UserID, err)
// 	}
// 	logger.Printf(ctx, "Manager Retrieved %v", hostelManager.ID)
// 	h.ManagerID = hostelManager.ID

// 	//Create the hostel
// 	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
// 		logger.Info("Creating hostel in transaction")
// 		hostel := Hostel{
// 			ManagerID:             h.ManagerID,
// 			Name:                  *h.Name,
// 			UniversityID:          h.UniversityID,
// 			Address:               *h.Address,
// 			City:                  *h.City,
// 			State:                 *h.State,
// 			Country:               *h.Country,
// 			Description:           *h.Description,
// 			NumberOfUnits:         *h.NumberOfUnits,
// 			NumberOfOccupiedUnits: *h.NumberOfOccupiedUnits,
// 			NumberOfBedrooms:      *h.NumberOfBedrooms,
// 			NumberOfBathrooms:     *h.NumberOfBathrooms,
// 			Kitchen:               *h.Kitchen,
// 			FloorSpace:            *h.FloorSpace,
// 		}

// 		//Create hostel together with image
// 		if err := tx.Model(&Hostel{}).Create(&hostel).Error; err != nil {
// 			logger.Error("Hostel creation failed")
// 			return err
// 		}
// 		logger.Printf(ctx, "Hostel created %v", hostel.ID)
// 		// Add the amenities

// 		amenitiesArr := make([]uint, len(h.Amenities))
// 		for _, ammenity := range h.Amenities {
// 			amenitiesArr = append(amenitiesArr, ammenity.ID)
// 		}
// 		amenities := []reference.ReferenceHostelAmenities{}
// 		if err := tx.Where("id IN ?", amenitiesArr).Find(&amenities).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Model(&hostel).Association("Amenities").Append(amenities); err != nil {
// 			return err
// 		}
// 		logger.Info("Amenities added")
// 		// Add the hostel images

// 		// Add the hostel fee
// 		breakdown := make(map[string]float64)
// 		for k, v := range h.HostelFee.Breakdown {
// 			breakdown[k] = v
// 		}

// 		//convert the map to json

// 		hostelFee := HostelFee{
// 			HostelID:    hostel.ID,
// 			TotalAmount: h.HostelFee.TotalAmount,
// 			PaymentPlan: h.HostelFee.PaymentPlan,
// 			Breakdown:   breakdown,
// 		}
// 		logger.Info(hostelFee.Breakdown)
// 		if tx.Model(&HostelFee{}).Create(&hostelFee).Error != nil {
// 			logger.Error("Hostel fee creation failed")
// 			return err
// 		}
// 		logger.Info("Hostel fee created")

// 		// In the main thread, receive file paths from the channel and create HostelImage records
// 		hostelImages := make([]HostelImage, len(h.Images))
// 		//logger.Info(filePathsChan)
// 		for _, path := range hostelImagesSlice {
// 			hostelImage := HostelImage{
// 				HostelID: h.Hostel.ID,
// 			}
// 			logger.Info("In transaction--->", path)
// 			hostelImage.ImageURL = path
// 			//logger.Info(hostelImage)
// 			hostelImages = append(hostelImages, hostelImage)
// 		}

// 		// Create HostelImage records
// 		logger.Info("Creating hostel images", hostelImages)
// 		if err := tx.Model(&HostelImage{}).Create(&hostelImages).Error; err != nil {
// 			return err
// 		}
// 		logger.Info("Hostel images created")

// 		h.Hostel = hostel

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	select {
// 	case err := <-errorChan:
// 		// Handle the error (e.g., log it)
// 		logger.Error(err)
// 		// Abort the request by returning an error response
// 		return err
// 	default:
// 		// No error, continue processing
// 		logger.Info("No error")}

// 	// Get the hostel again so that the related fields a preloaded
// 	logger.Info("Getting hostel")
// 	err = db.Model(&Hostel{}).Preload(clause.Associations).Where("id = ?", h.Hostel.ID).First(&h.Hostel).Error
// 	if err != nil {
// 		return err
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Create creates a new hostel
func (h *HostelSerializer) Create(ctx *gin.Context, db *gorm.DB, session *session.Session) error {
	logger := common.NewLogger()

	//Validate the fields
	if err := h.Validate(); err != nil {
		return err
	}
	// Get the manager id from the auth payload
	hostelManager, err := getManager(ctx, db)
	if err != nil {
		return err
	}

	h.ManagerID = hostelManager.ID

	//Process the uploaded images

	filePaths, err := utils.ProcessUploadedImages(h.Images, session)
	if err != nil {
		logger.Info("Image processing failed")
		return err
	}

	//Create the hostel
	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		logger.Info("Creating hostel in transaction")
		hostel := Hostel{
			ManagerID:             h.ManagerID,
			Name:                  *h.Name,
			UniversityID:          h.UniversityID,
			Address:               *h.Address,
			City:                  *h.City,
			State:                 *h.State,
			Country:               *h.Country,
			Description:           *h.Description,
			NumberOfUnits:         *h.NumberOfUnits,
			NumberOfOccupiedUnits: *h.NumberOfOccupiedUnits,
			NumberOfBedrooms:      *h.NumberOfBedrooms,
			NumberOfBathrooms:     *h.NumberOfBathrooms,
			Kitchen:               *h.Kitchen,
			FloorSpace:            *h.FloorSpace,
		}
		//Create hostel together with image
		if err := tx.Model(&Hostel{}).Create(&hostel).Error; err != nil {
			logger.Error("Hostel creation failed")
			return err
		}
		logger.Printf(ctx, "Hostel created %v", hostel.ID)
		// Add the amenities

		amenitiesArr := make([]uint, len(h.Amenities))
		for _, ammenity := range h.Amenities {
			amenitiesArr = append(amenitiesArr, ammenity.ID)
		}
		amenities := []reference.ReferenceHostelAmenities{}
		if err := tx.Where("id IN ?", amenitiesArr).Find(&amenities).Error; err != nil {
			return err
		}

		if err := tx.Model(&hostel).Association("Amenities").Append(amenities); err != nil {
			return err
		}
		logger.Info("Amenities added")
		// Add the hostel images
		hostelImages := make([]HostelImage, len(filePaths))

		for i, image := range filePaths {
			hostelImages[i] = HostelImage{}
			hostelImages[i].ImageURL = image
			hostelImages[i].HostelID = hostel.ID
		}

		if err := tx.Model((&HostelImage{})).Create(&hostelImages).Error; err != nil {
			return err
		}
		// Add the hostel fee
		breakdown := make(map[string]float64)
		for k, v := range h.HostelFee.Breakdown {
			breakdown[k] = v
		}

		//convert the map to json

		hostelFee := HostelFee{
			HostelID:    hostel.ID,
			TotalAmount: h.HostelFee.TotalAmount,
			PaymentPlan: h.HostelFee.PaymentPlan,
			Breakdown:   breakdown,
		}
		logger.Info(hostelFee.Breakdown)
		if tx.Model(&HostelFee{}).Create(&hostelFee).Error != nil {
			logger.Error("Hostel fee creation failed")
			return err
		}
		logger.Info("Hostel fee created")

		h.Hostel = hostel
		return nil
	})
	if err != nil {
		return err
	}
	// Get the hostel again so that the related fields a preloaded

	err = db.Model(&Hostel{}).Preload(clause.Associations).Where("id = ?", h.Hostel.ID).First(&h.Hostel).Error
	if err != nil {
		return err
	}

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

	// Register the custom validation function
	validate.RegisterValidation("ValidateImageExtension", ValidateImageExtension)

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
func (h *HostelSerializer) Response() map[string]interface{} {
	amenities := make([]map[string]interface{}, len(h.Hostel.Amenities))
	for i, ammenity := range h.Hostel.Amenities {
		amenities[i] = map[string]interface{}{
			"id":   ammenity.ID,
			"name": ammenity.Name,
		}
	}
	imageURLs := make([]string, len(h.Hostel.HostelImages))
	for i, image := range h.Hostel.HostelImages {
		imageURLs[i] = image.ImageURL
	}

	hostelFee := map[string]interface{}{
		"total_amount": h.Hostel.HostelFee.TotalAmount,
		"payment_plan": h.Hostel.HostelFee.PaymentPlan,
		"breakdown":    h.Hostel.HostelFee.Breakdown,
	}

	return map[string]interface{}{
		"uid":                      h.Hostel.UID,
		"created_at":               h.Hostel.CreatedAt,
		"updated_at":               h.Hostel.UpdatedAt,
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

func (h *HostelSerializer) Get(db *gorm.DB, hostelUID uuid.UUID) error {
	err := db.Model(&Hostel{}).Preload(clause.Associations).Where("uid = ?", hostelUID).First(&h.Hostel).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *HostelSerializer) List(db *gorm.DB, queryParams HostelQueryParams) ([]Hostel, error) {
	hostels, err := h.Filter(db, queryParams)
	if err != nil {
		return nil, err
	}

	return hostels, nil
}

func (h HostelSerializer) ResponseMany(hostels []Hostel) []map[string]interface{} {
	hostelsResponse := make([]map[string]interface{}, len(hostels))
	for i, hostel := range hostels {
		serializer := HostelSerializer{Hostel: hostel}
		hostelsResponse[i] = serializer.Response()
	}
	return hostelsResponse
}

func (h *HostelSerializer) Update(ctx *gin.Context, db *gorm.DB, session *session.Session, hostelUID uuid.UUID) error {
	logger := common.NewLogger()
	var hostel Hostel

	if err := h.Validate(); err != nil {
		return err
	}
	// Get the manager id from the auth payload
	hostelManager, err := getManager(ctx, db)
	if err != nil {
		return err
	}

	updatedFields := h.getUpdatedFields()
	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		err := tx.Model(&Hostel{}).Preload(clause.Associations).Where("uid = ?", hostelUID).First(&hostel).Error
		if err != nil {
			return err
		}
		if hostel.ManagerID != hostelManager.ID {
			return fmt.Errorf("hostel does not belong to manager")
		}

		if err := tx.Model(&hostel).Updates(updatedFields).Error; err != nil {
			return err
		}
		if h.Amenities != nil {
			logger.Info("Updating Amenities")
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
	h.Hostel = hostel
	if err != nil {
		return err
	}
	return nil
}

func (h *HostelSerializer) updateAmenities(hostel *Hostel, tx *gorm.DB) ([]uint, error) {
	// Create a map to track amenities by ID for faster lookup
	logger := common.NewLogger()
	existingAmenities := make(map[uint]bool)
	for _, amenity := range hostel.Amenities {
		existingAmenities[amenity.ID] = true
	}

	// Create a map to track new amenities by ID and collect their IDs
	newAmenities := make(map[uint]bool)
	amenitiesArr := make([]uint, 0)

	for _, newAmenity := range h.Amenities {
		if !existingAmenities[newAmenity.ID] {
			amenitiesArr = append(amenitiesArr, newAmenity.ID)
		}
		newAmenities[newAmenity.ID] = true
	}

	// Remove amenities that are not in the new list
	for _, amenity := range hostel.Amenities {
		if !newAmenities[amenity.ID] {
			if err := tx.Model(&hostel).Association("Amenities").Delete(&amenity); err != nil {
				return nil, err
			}
		}
	}
	logger.Info(amenitiesArr)

	return amenitiesArr, nil
}

// Function to filter hostels based on query parameters
func (h *HostelSerializer) Filter(db *gorm.DB, queryParams HostelQueryParams) ([]Hostel, error) {
	// Initialize the query with the Hostel model
	var hostels []Hostel

	query := db.Model(&Hostel{}).Preload(clause.Associations)

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

	// You can also apply other filters such as name, address, etc., as needed
	err := query.Preload(clause.Associations).Find(&hostels).Error
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
