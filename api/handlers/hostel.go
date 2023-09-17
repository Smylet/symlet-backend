package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

// @Summary Create a hostel
// @Description Create a new hostel
// @Tags Hostels
// @Accept multipart/form-data
// @Produce json
// @Param hostel body hostel.HostelSerializer true "Hostel object to create"
// @Param hostel_images formData multipart.FileHeader false "Hostel images"
// @Success 201 {object} hostel.HostelSerializer "Created hostel"
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /hostels [post]
func (server *Server) CreateHostel(c *gin.Context) {
	var HostelSerializer hostel.HostelSerializer
	//print header
	file, err := c.FormFile("hostel_images")
	if err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid hostel images")
		return
	}
	errs := utils.CustomBinder(c, &HostelSerializer)
	if errs != nil {
		utils.RespondWithError(c, http.StatusBadRequest, errs.Error(), "Invalid hostel data")
		return
	}

	HostelSerializer.Images = append(HostelSerializer.Images, file)
	//Get uploaded file
	err = HostelSerializer.CreateTx(c, server.db, server.session)

	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create hostel")
		return
	}

	utils.RespondWithSuccess(c, 201, HostelSerializer.Response(), "Email confirmed")
}

// @Summary Get a hostel
// @Description Get a hostel by uid
// @Tags Hostels
// @Accept json
// @Produce json
// @Param uid path string true "Hostel uid"
// @Success 200 {object} hostel.HostelSerializer "Hostel"
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /hostels/{uid} [get]
func (server *Server) GetHostel(c *gin.Context) {
	var HostelSerializer hostel.HostelSerializer
	uidString := c.Param("uid")
	if uidString == "" {
		utils.RespondWithError(c, 400, "Hostel uid is required", "")
		return
	}
	uid, err := uuid.Parse(uidString)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid hostel uid")
		return
	}
	err = HostelSerializer.GetHostelTx(server.db, uid)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to get hostel")
		return
	}
	utils.RespondWithSuccess(c, 200, HostelSerializer.Response(), "Hostel retrieved successfully")
}

// @Summary List hostels with optional filters
// @Description Get a list of hostels with optional filters
// @Tags Hostels
// @Accept json
// @Produce json
// @Param name query string false "Hostel name"
// @Param description query string false "Description"
// @Param university_id query uint false "University ID"
// @Param address query string false "Address"
// @Param city query string false "City"
// @Param state query string false "State"
// @Param country query string false "Country"
// @Param manager_id query uint false "Manager ID"
// @Param number_of_units query uint false "Number of units"
// @Param number_of_bedrooms query uint false "Number of bedrooms"
// @Param number_of_bathrooms query uint false "Number of bathrooms"
// @Param kitchen query string false "Kitchen type (shared, none, private)"
// @Param floor_space query uint false "Floor space"
// @Param hostel_fee_total_min query number false "Minimum hostel fee total"
// @Param hostel_fee_total_max query number false "Maximum hostel fee total"
// @Param hostel_fee_plan query string false "Hostel fee plan"
// @Param has_amenities query bool false "Has amenities"
// @Param security_rating_min query number false "Minimum security rating"
// @Param security_rating_max query number false "Maximum security rating"
// @Param location_rating_min query number false "Minimum location rating"
// @Param location_rating_max query number false "Maximum location rating"
// @Param general_rating_min query number false "Minimum general rating"
// @Param general_rating_max query number false "Maximum general rating"
// @Param amenities_rating_min query number false "Minimum amenities rating"
// @Param amenities_rating_max query number false "Maximum amenities rating"
// @Param water_rating_min query number false "Minimum water rating"
// @Param water_rating_max query number false "Maximum water rating"
// @Param electricity_rating_min query number false "Minimum electricity rating"
// @Param electricity_rating_max query number false "Maximum electricity rating"
// @Param caretaker_rating_min query number false "Minimum caretaker rating"
// @Param caretaker_rating_max query number false "Maximum caretaker rating"
// @Success 200 {object} hostel.HostelSerializer "Hostels"
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /hostels [get]
func (server *Server) ListHostels(c *gin.Context) {
	var HostelSerializer hostel.HostelSerializer
	var queryParams hostel.HostelQueryParams
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid query parameters")
		return
	}

	hostels, err := HostelSerializer.ListHostelsTx(server.db, queryParams)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error(), "Failed to get hostels")
		return
	}
	if len(hostels) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, nil, "No hostels found")
		return
	}
	hostelsResponse := HostelSerializer.ResponseMany(hostels)
	utils.RespondWithSuccess(c, http.StatusOK, hostelsResponse, "Hostels retrieved successfully")
}

// @Summary Update a hostel
// @Description Update a hostel by uid
// @Tags Hostels
// @Accept json
// @Produce json
// @Param uid path string true "Hostel uid"
// @Param hostel body hostel.HostelSerializer true "Hostel object to update"
// @Success 200 {object} hostel.HostelSerializer "Updated hostel"
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /hostels/{uid} [put]
func (server *Server) UpdateHostel(c *gin.Context) {
	var HostelSerializer hostel.HostelSerializer
	uidString := c.Param("uid")
	errs := utils.CustomBinder(c, &HostelSerializer)
	if errs != nil {
		utils.RespondWithError(c, http.StatusBadRequest, errs.Error(), "Invalid hostel data")
		return
	}

	uid, err := uuid.Parse(uidString)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid hostel uid")
		return
	}

	err = HostelSerializer.UpdateHostelTx(c, server.db, server.session, uid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error(), "Failed to update hostel")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, HostelSerializer.Response(), "Hostel updated successfully")
}

// @Summary Delete a hostel
// @Description Delete a hostel by uid
// @Tags Hostels
// @Accept json
// @Produce json
// @Param uid path string true "Hostel uid"
// @Success 204 "No content"
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /hostels/{uid} [delete]
func (server *Server) DeleteHostel(c *gin.Context) {
	var hostel hostel.Hostel
	uidString := c.Param("uid")
	uid, err := uuid.Parse(uidString)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid hostel uid")
		return
	}

	err = server.db.Where("uid = ?", uid).First(&hostel).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error(), "Failed to get hostel")
		return
	}

	err = server.db.Delete(&hostel).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error(), "Failed to delete hostel")
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
	utils.RespondWithSuccess(c, http.StatusNoContent, nil, "Hostel deleted successfully")
}
