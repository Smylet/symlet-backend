package users

import (
	"github.com/Smylet/symlet-backend/api/users/routers"
	"github.com/gin-gonic/gin"
)

// Register all the routes

func RegisterRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", routers.Register)
		userRoutes.POST("/login", routers.Login)
		userRoutes.POST("/logout", routers.Logout)
		userRoutes.POST("/confirm-email", routers.ConfirmEmail)
		userRoutes.POST("/email/confirmation/resend", routers.ResendEmailConfirmation)
		userRoutes.POST("/password-reset", routers.RequestPasswordReset)
		userRoutes.PUT("/password-change", routers.ChangePassword)

		// USER PROFILE MANAGEMENT ENDPOINTS
		userRoutes.GET("/:username/profile", routers.GetUserProfile)
		userRoutes.PUT("/:username/profile", routers.EditUserProfile)
		userRoutes.DELETE("/:username", routers.DeleteUserProfile)
		userRoutes.GET("/:username/history", routers.ViewProfileEditHistory)
		userRoutes.POST("/:username/backup", routers.BackupUserProfile)
		userRoutes.GET("/:username/backups", routers.ListProfileBackups)
		userRoutes.PUT("/:username/restore", routers.RestoreUserProfile)
		userRoutes.GET("/:username/export", routers.ExportUserProfile)
		userRoutes.PUT("/:username/deactivate", routers.DeactivateAccount)
		userRoutes.PUT("/:username/reactivate", routers.ReactivateAccount)

		// USER PROFILE PICTURE ENDPOINTS
		userRoutes.POST("/:username/picture", routers.UploadProfilePicture)
		userRoutes.PUT("/:username/picture", routers.UpdateProfilePicture)
		userRoutes.DELETE("/:username/picture", routers.DeleteProfilePicture)

		// USER NOTIFICATION ENDPOINTS
		userRoutes.GET("/:username/notifications", routers.GetNotification)
		userRoutes.PUT("/:username/notifications/notificationId/read", routers.UpdateNotification)
		userRoutes.PUT("/:username/notifications/settings", routers.UpdateNotificationSettings)
		userRoutes.GET("/:username/notifications/settings", routers.GetNotificationSettings)

		// USER SEARCH ENDPOINTS
		userRoutes.GET("/search", routers.SearchUsers)

		// USER PRIVACY SETTINGS ENDPOINTS
		userRoutes.GET("/:username/privacy", routers.GetPrivacySettings)
		userRoutes.PUT("/:username/privacy", routers.UpdatePrivacySettings)

	}

	userRoute := r.Group("/user")
	{
		// USER 2FA ENDPOINTS
		userRoute.POST("/2fa/setup", routers.Setup2FA)
		userRoute.POST("/2fa/verify", routers.Verify2FA)
	}

	profileRoutes := r.Group("/profiles")
	{
		ProfileRegister(profileRoutes) // Profile operations
	}
}

// Profile-related routes
func ProfileRegister(r *gin.RouterGroup) {
	r.GET("/:username", routers.GetProfile)
	// Add other routes as needed
}
