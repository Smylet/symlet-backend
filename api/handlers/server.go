package handlers

import (
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	config utils.Config
	cron   *cron.Cron
	db     *gorm.DB
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config, db *gorm.DB) (*Server, error) {

	server := &Server{
		config: config,
		cron:   cron.New(),
		db:     db,
	}

	server.registerUserRoutes()
	server.registerReviewRoutes()
	return server, nil
}

func (server *Server) registerUserRoutes() {

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", server.Register)
		userRoutes.POST("/login", server.Login)
		userRoutes.POST("/logout", server.Logout)
		userRoutes.POST("/confirm-email", server.ConfirmEmail)
		userRoutes.POST("/email/confirmation/resend", server.ResendEmailConfirmation)
		userRoutes.POST("/password-reset", server.RequestPasswordReset)
		userRoutes.PUT("/password-change", server.ChangePassword)

		// USER PROFILE MANAGEMENT ENDPOINTS
		userRoutes.GET("/:username/profile", server.GetUserProfile)
		userRoutes.PUT("/:username/profile", server.EditUserProfile)
		userRoutes.DELETE("/:username", server.DeleteUserProfile)
		userRoutes.GET("/:username/history", server.ViewProfileEditHistory)
		userRoutes.GET("/:username/backups", server.ListProfileBackups)
		userRoutes.PUT("/:username/restore", server.RestoreUserProfile)
		userRoutes.GET("/:username/export", server.ExportUserProfile)
		userRoutes.PUT("/:username/deactivate", server.DeactivateAccount)
		userRoutes.PUT("/:username/reactivate", server.ReactivateAccount)

		// USER PROFILE PICTURE ENDPOINTS
		userRoutes.POST("/:username/picture", server.UploadProfilePicture)
		userRoutes.PUT("/:username/picture", server.UpdateProfilePicture)
		userRoutes.DELETE("/:username/picture", server.DeleteProfilePicture)

		// USER NOTIFICATION ENDPOINTS
		userRoutes.GET("/:username/notifications", server.GetNotification)
		userRoutes.PUT("/:username/notifications/notificationId/read", server.UpdateNotification)
		userRoutes.PUT("/:username/notifications/settings", server.UpdateNotificationSettings)
		userRoutes.GET("/:username/notifications/settings", server.GetNotificationSettings)

		// USER SEARCH ENDPOINTS
		userRoutes.GET("/search", server.SearchUsers)

		// USER PRIVACY SETTINGS ENDPOINTS
		userRoutes.GET("/:username/privacy", server.GetPrivacySettings)
		userRoutes.PUT("/:username/privacy", server.UpdatePrivacySettings)

	}

	userRoute := r.Group("/user")
	{
		// USER 2FA ENDPOINTS
		userRoute.POST("/2fa/setup", server.Setup2FA)
		userRoute.POST("/2fa/verify", server.Verify2FA)
	}

	profileRoutes := r.Group("/profiles")
	{
		profileRoutes.GET("/:username", server.GetProfile)
	}

	server.router = r
}

func (server *Server) registerReviewRoutes() {

}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
