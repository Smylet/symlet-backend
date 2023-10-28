package handlers

import (
	"net/http"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetMe(c *gin.Context) {
	var userSerializer users.UserSerializer

	payload, err := users.GetAuthPayloadFromCtx(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, err.Error(), "Unauthorized")
	}

	err = userSerializer.FindCurrentUser(c, server.db, payload.UserID)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to get user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioGet), "User retrieved successfully")
}

func (server *Server) GetUser(c *gin.Context) {
	var userSerializer users.UserSerializer

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	_, err := userSerializer.FindByUID(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to get user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioGet), "User retrieved successfully")
}

func (server *Server) SearchUsers(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioSearch, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	userList, err := userSerializer.Search(c, server.db)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error(), "Failed to search users")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.ListResponse(userList), "Users retrieved successfully")
}

func (server *Server) GetUsers(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioList, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	userList, err := userSerializer.List(c, server.db)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to get users")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.ListResponse(userList), "Users retrieved successfully")
}

func (server *Server) GetPreferences(c *gin.Context) {
	var userSerializer users.UserSerializer

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	preferences, err := userSerializer.GetPreferences(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to get users")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, preferences, "Preferences retrieved successfully")
}

func (server *Server) UpdatePreferences(c *gin.Context) {

	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioUpdatePreferences, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	err = userSerializer.UpdatePreferences(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to update preferences")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "", "Preferences updated successfully")

}

func (server *Server) DeletePreferences(c *gin.Context) {

	var userSerializer users.UserSerializer

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	err := userSerializer.DeletePreferences(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to delete preferences")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "", "Preferences deleted successfully")
}
func (server *Server) GetPastSearches(c *gin.Context) {

	var userSerializer users.UserSerializer

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	pastSearches, err := userSerializer.GetPastSearches(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to get past searches")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, pastSearches, "Past searches retrieved successfully")
}

func (server *Server) AddPastSearch(c *gin.Context) {

	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioAddPastSearch, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	err = userSerializer.AddPastSearch(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to add past search")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioAddPastSearch), "Past search added successfully")
}

func (server *Server) ClearPastSearches(c *gin.Context) {

	var userSerializer users.UserSerializer

	uid := c.Param("uid")
	if uid == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID", "Please provide a valid user ID")
		return
	}

	err := userSerializer.ClearPastSearches(c, server.db, uid)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to clear past searches")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "", "Past searches cleared successfully")
}
