// USER PROFILE MANAGEMENT ENDPOINTS
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

// USER PROFILE MANAGEMENT ENDPOINTS

// @Summary Create a new user profile
// @Description Create a new user profile
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param profile body users.ProfileSerializer true "Profile object to create"
// @Success 201 {object} users.ProfileSerializer
// @Failure 401 {object} utils.ErrorMessage
// @Failure 400 {object} utils.ErrorMessage "Bad request"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /users/{uid}/profile [post]
func (server *Server)CreateUserProfile(c *gin.Context){
	var profileSerializer users.ProfileSerializer
	fmt.Print("Creating profile")
	//print header
	file, err := c.FormFile("image")
	if err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid profile image")
		return
	}
	errs := common.CustomBinder(c, &profileSerializer)
	if errs != nil {
		utils.RespondWithError(c, http.StatusBadRequest, errs.Error(), "Invalid profile data")
		return
	}
	profileSerializer.Image = file

	err = profileSerializer.Create(c, server.db, server.session)

	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create profile")
		return
	}

	utils.RespondWithSuccess(c, 201, profileSerializer.Response(), "Profile created successfully")
}


func (server *Server) GetUserProfile(c *gin.Context) {
	// Handle logic to retrieve the user profile
	var profileSerializer users.ProfileSerializer
	uidString := c.Param("uid")
	if uidString == "" {
		utils.RespondWithError(c, 400, "profile uid is required", "")
		return
	}

	err := profileSerializer.Get(c, server.db, uidString)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to retrieve profile")
		return
	}
	utils.RespondWithSuccess(c, 200, profileSerializer.Response(), "Profile retrieved successfully")

}

func (server *Server) EditUserProfile(c *gin.Context) {
	// Handle logic to edit the user profile
	c.JSON(200, gin.H{
		"message": "User profile edited successfully",
	})
}

func (server *Server) DeleteUserProfile(c *gin.Context) {
	// Handle logic to delete the user profile
	c.JSON(200, gin.H{
		"message": "User profile deleted successfully",
	})
}

func (server *Server) ViewProfileEditHistory(c *gin.Context) {
	// Handle logic to view the user profile's edit history
	c.JSON(200, gin.H{
		"message": "User profile edit history retrieved",
	})
}

func (server *Server) BackupUserProfile(c *gin.Context) {
	// Handle logic to create a backup of the user's profile
}

func (server *Server) ListProfileBackups(c *gin.Context) {
	// Handle logic to list all backups of the user's profile
	c.JSON(200, gin.H{
		"message": "Profile backups listed",
	})
}

func (server *Server) RestoreUserProfile(c *gin.Context) {
	// Handle logic to restore the user's profile from a backup
	c.JSON(200, gin.H{
		"message": "User profile restored successfully",
	})
}

func (server *Server) ExportUserProfile(c *gin.Context) {
	// Handle logic to export user's profile data
	c.JSON(200, gin.H{
		"message": "User profile data exported",
	})
}

func (server *Server) DeactivateAccount(c *gin.Context) {
	// Handle logic to deactivate user's account
	c.JSON(200, gin.H{
		"message": "User account deactivated",
	})
}

func (server *Server) ReactivateAccount(c *gin.Context) {
	// Handle logic to reactivate a user's deactivated account
	c.JSON(200, gin.H{
		"message": "User account reactivated",
	})
}

// USER PROFILE PICTURE ENDPOINTS

func (server *Server) UploadProfilePicture(c *gin.Context) {
	// Handle logic to upload a new profile picture
	c.JSON(200, gin.H{
		"message": "Profile picture uploaded successfully",
	})
}

func (server *Server) UpdateProfilePicture(c *gin.Context) {
	// Handle logic to update user's profile picture
	c.JSON(200, gin.H{
		"message": "Profile picture updated",
	})
}

func (server *Server) DeleteProfilePicture(c *gin.Context) {
	// Handle logic to delete user's profile picture
	c.JSON(200, gin.H{
		"message": "Profile picture deleted",
	})
}

func (server *Server) SearchUsers(c *gin.Context) {
	// ...
}

func (server *Server) GetPrivacySettings(c *gin.Context) {
	// ...
}

func (server *Server) UpdatePrivacySettings(c *gin.Context) {
	// ...
}
