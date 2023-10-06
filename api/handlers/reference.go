package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

// @Summary List amenities
// @Description List all amenities
// @Tags Amenities
// @Produce json
// @Success 200 {object} reference.AmenitySerializer
// @Router /references/amenities [get]
func (server *Server) ListAmenities(c *gin.Context) {
	var amenitySerializer reference.AmenitySerializer

	amenities, err := amenitySerializer.List(server.db)

	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to list amenities")
		return
	}
	response := amenitySerializer.ResponseMany(amenities)

	utils.RespondWithSuccess(c, 200, response, "Amenities listed successfully")
}


// @Summary List universities
// @Description List all universities
// @Tags Universities
// @Produce json
// @Param name query string false "University name"
// @Param code query string false "University code"
// @Param city query string false "University city"
// @Param country query string false "University country"
// @Param state query string false "University state"
// @Success 200 {object} reference.UniversitySerializer
// @Router /references/universities [get]
func (server *Server) ListUniversities(c *gin.Context) {
	var universitySerializer reference.UniversitySerializer
	var queryParams reference.UniversityQueryParams

	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid query parameters")
		return
	}

	universities, err := universitySerializer.List(server.db, queryParams)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to list universities")
		return
	}

	response := universitySerializer.ResponseMany(universities)
	if len(response) == 0 {
		utils.RespondWithSuccess(c, 200, response, "No universities found")
		return
	}
	utils.RespondWithSuccess(c, 200, response, "Universities listed successfully")
}