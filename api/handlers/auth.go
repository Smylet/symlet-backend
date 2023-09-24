// USER ACCOUNT & AUTHENTICATION ENDPOINTS
package handlers

import (
	"log"
	"time"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
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
func (server *Server) Register(c *gin.Context) {
	var req users.CreateUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid request body")
		return
	}

	if status := users.ValidateRegisterUserReq(req); !status.Valid {
		utils.RespondWithError(c, 400, "", status.Message)
		return
	}

	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to hash password")
		return
	}

	arg := users.CreateUserTxParams{
		CreateUserReq: users.CreateUserReq{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		},
		AfterCreate: func(user users.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.task.DistributeTaskSendVerifyEmail(c, taskPayload, opts...)
		},
	}
	userRepo := users.NewUserRepository(server.db)

	txResult, err := userRepo.CreateUserTx(c, arg)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create user")
		return
	}

	user := users.UserSerializer{User: txResult.User}
	utils.RespondWithSuccess(c, 200, user.Response(), "User created successfully")
}

// ConfirmEmail godoc
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
func (server *Server) ConfirmEmail(c *gin.Context) {
	var req users.ConfirmVerifyEmailParams

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid request body")
		return
	}

	userRepo := users.NewUserRepository(server.db)

	if err := userRepo.VerifyEmailTx(c, users.ConfirmVerifyEmailParams{
		UserID:     req.UserID,
		VerEmailID: req.VerEmailID,
		SecretCode: req.SecretCode,
	}); err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to verify email")
		return
	}

	utils.RespondWithSuccess(c, 200, nil, "Email confirmed")
}

// Login godoc
// @Summary Login a user
// @Description Login a user with the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Param req body users.LoginReq true "User login information"
// @Success 200 {object} utils.SuccessMessage
// @Failure 400 {object} utils.ErrorMessage
// @Failure 500 {object} utils.ErrorMessage
// @Router /users/login [post]
func (server *Server) Login(c *gin.Context) {
	var req users.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid request body")
		return
	}

	userRepo := users.NewUserRepository(server.db)
	user, err := userRepo.FindUser(c, users.FindUserParams{
		User: users.User{
			Email: req.Email,
		},
	})
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to find user")
		return
	}

	err = common.CheckPassword(req.Password, user.Password)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid password")
		return
	}

	accessToken, accessPayload, err := server.token.CreateToken(
		user.Username,
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create access token")
		return
	}

	refreshToken, refreshPayload, err := server.token.CreateToken(
		user.Username,
		user.ID,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create refresh token")
		return
	}

	session, err := userRepo.CreateSession(c, users.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIP:     c.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create session")
		return
	}

	response := users.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  user,
	}
	utils.RespondWithSuccess(c, 200, response, "Logged in successfully")
}

func (server *Server) Logout(c *gin.Context) {
	// Handle user logout logic
	// This would typically involve revoking tokens or clearing sessions.
	c.JSON(200, gin.H{
		"message": "Logged out successfully",
	})
}

func (server *Server) ResendEmailConfirmation(c *gin.Context) {
	// Handle resending of email confirmation logic
	c.JSON(200, gin.H{
		"message": "Email confirmation resent",
	})
}

func (server *Server) RequestPasswordReset(c *gin.Context) {
	// Handle password reset request logic
	c.JSON(200, gin.H{
		"message": "Password reset link sent",
	})
}

func (server *Server) ChangePassword(c *gin.Context) {
	authPayload := c.MustGet(users.AuthorizationPayloadKey).(*token.Payload)
	log.Println(authPayload)
	// Handle password change logic
	c.JSON(200, gin.H{
		"message": "Password changed successfully",
	})
}

func (server *Server) Setup2FA(c *gin.Context) {
	// Handle 2FA setup logic
	c.JSON(200, gin.H{
		"message": "2FA setup successfully",
	})
}

func (server *Server) Verify2FA(c *gin.Context) {
	// Handle 2FA verification logic
	c.JSON(200, gin.H{
		"message": "2FA verified",
	})
}

func (server *Server) GetProfile(c *gin.Context) {
	username := c.Param("username")
	// Fetch the profile based on the username
	// This would involve querying your data source for profile information
	c.JSON(200, gin.H{
		"username": username,
		"bio":      "Sample bio for the user",
	})
}
