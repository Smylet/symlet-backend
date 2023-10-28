// USER ACCOUNT & AUTHENTICATION ENDPOINTS
package handlers

import (
	"net/http"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Param req body users.CreateUserReq true "User registration information"
// @Success 200 {object} utils.SuccessMessage "User created successfully"
// @Failure 400 {object} utils.ErrorMessage "Invalid request body or validation failure"
// @Failure 500 {object} utils.ErrorMessage "Server error or unexpected issues"
// @Router /users/register [post]
func (server *Server) CreateUser(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioCreate, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	if err := userSerializer.Create(c, server.db, server.task, server.cron); err != nil {
		if err != nil {
			utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to create user")
			return
		}

	}

	utils.RespondWithSuccess(c, http.StatusCreated, userSerializer.Response(common.ScenarioCreate), "User created successfully")
}

// verifyEmail godoc
// @Summary Confirm a user's email
// @Description Confirm a user's email using verification parameters.
// @Tags users
// @Accept  json
// @Produce  json
// @Param userID query string true "User ID for email verification"
// @Param verEmailID query string true "Verification Email ID"
// @Param secretCode query string true "Secret verification code"
// @Success 200 {object} utils.SuccessMessage "Email successfully confirmed"
// @Failure 400 {object} utils.ErrorMessage "Invalid request or parameters"
// @Failure 500 {object} utils.ErrorMessage "Server error or unexpected issues"
// @Router /users/confirm-email [post]
func (server *Server) verifyEmail(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioVerifyEmail, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	if err := userSerializer.VerifyEmail(c, server.db); err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to confirm email")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioVerifyEmail), "Email successfully confirmed")
}

// LoginUser godoc
// @Summary LoginUser a user
// @Description LoginUser a user with the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Param req body users.LoginReq true "User login information"
// @Success 200 {object} utils.SuccessMessage
// @Failure 400 {object} utils.ErrorMessage
// @Failure 500 {object} utils.ErrorMessage
// @Router /users/login [post]
func (server *Server) LoginUser(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioLogin, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.LoginUser(c, server.db, server.token, server.config, server.redisClient)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to login")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioLogin), "Logged in successfully")
}

func (server *Server) RenewAccessToken(c *gin.Context) {

	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioRenewAccessToken, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.RenewAccessToken(c, server.db, server.token, server.config)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to renew access token")
		return
	}

	utils.RespondWithSuccess(c, 200, userSerializer.Response(common.ScenarioRenewAccessToken), "Access token renewed successfully")
}

func (server *Server) ForgotPassword(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioForgotPassword, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.ForgotPassword(c, server.db, server.task, server.token)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to request password reset")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioForgotPassword), "Password reset link sent successfully")
}

func (server *Server) PasswordReset(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioPasswordReset, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.PasswordReset(c, server.db, server.token)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to reset password")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioPasswordReset), "Password reset successfully")

}
func (server *Server) LogoutUser(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioLogout, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.LogoutUser(c, server.db, server.redisClient)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to logout")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioLogout), "Logged out successfully")
}

func (server *Server) ResendEmailVerification(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioResendEmailConfirmation, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.ResendEmailVerification(c, server.db, server.task)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to resend email confirmation")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioResendEmailConfirmation), "Email confirmation resent successfully")

}

func (server *Server) ChangePassword(c *gin.Context) {
	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioLogin, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.ChangePassword(c, server.db, server.task)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to change password")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioChangePassword), "Password changed successfully")
}

func (server *Server) Setup2FA(c *gin.Context) {

	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioSetup2FA, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.Setup2FA(c, server.db)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to setup 2FA")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioSetup2FA), "2FA setup successfully")
}

func (server *Server) Verify2FA(c *gin.Context) {

	var userSerializer users.UserSerializer

	err := common.CustomBinder(c, common.ScenarioVerify2FA, &userSerializer)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	err = userSerializer.Verify2FA(c, server.db, server.redisClient, server.token, server.config)
	if err != nil {
		utils.RespondWithError(c, userSerializer.StatusCode, err.Error(), "Failed to setup 2FA")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, userSerializer.Response(common.ScenarioVerify2FA), "2FA setup successfully")
}
