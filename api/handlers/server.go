package handlers

import (
	"github.com/Smylet/symlet-backend/api/users"
	_ "github.com/Smylet/symlet-backend/docs"

	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	router  *gin.Engine
	config  utils.Config
	cron    *cron.Cron
	db      *gorm.DB
	task    worker.TaskDistributor
	mailer  mail.EmailSender
	session *session.Session
	token   token.Maker
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config, db *gorm.DB, task worker.TaskDistributor, mailer mail.EmailSender, session *session.Session) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:  config,
		cron:    cron.New(),
		db:      db,
		task:    task,
		mailer:  mailer,
		session: session,
		token:   tokenMaker,
	}

	gin.SetMode(gin.ReleaseMode)

	server.registerRoutes()

	return server, nil
}

func (server *Server) registerRoutes() {
	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", server.Register)
		userRoutes.POST("/login", server.Login)
		userRoutes.POST("/logout", server.Logout)
		userRoutes.POST("/confirm-email", server.ConfirmEmail)
		userRoutes.POST("/email/confirmation/resend", server.ResendEmailConfirmation)
		userRoutes.POST("/password-reset", server.RequestPasswordReset)
		userRoutes.PUT("/password-change", users.AuthMiddleware(server.token), server.ChangePassword)

		// USER PROFILE MANAGEMENT ENDPOINTS
		userRoutes.GET("/:username/profile", users.AuthMiddleware(server.token), server.GetUserProfile)
		userRoutes.PUT("/:username/profile", users.AuthMiddleware(server.token), server.EditUserProfile)
		userRoutes.DELETE("/:username", users.AuthMiddleware(server.token), server.DeleteUserProfile)
		userRoutes.GET("/:username/history", users.AuthMiddleware(server.token), server.ViewProfileEditHistory)
		userRoutes.GET("/:username/backups", users.AuthMiddleware(server.token), server.ListProfileBackups)
		userRoutes.PUT("/:username/restore", users.AuthMiddleware(server.token), server.RestoreUserProfile)
		userRoutes.GET("/:username/export", users.AuthMiddleware(server.token), server.ExportUserProfile)
		userRoutes.PUT("/:username/deactivate", users.AuthMiddleware(server.token), server.DeactivateAccount)
		userRoutes.PUT("/:username/reactivate", users.AuthMiddleware(server.token), server.ReactivateAccount)

		// USER PROFILE PICTURE ENDPOINTS
		userRoutes.POST("/:username/picture", users.AuthMiddleware(server.token), server.UploadProfilePicture)
		userRoutes.PUT("/:username/picture", users.AuthMiddleware(server.token), server.UpdateProfilePicture)
		userRoutes.DELETE("/:username/picture", users.AuthMiddleware(server.token), server.DeleteProfilePicture)

		// USER NOTIFICATION ENDPOINTS
		userRoutes.GET("/:username/notifications", users.AuthMiddleware(server.token), server.GetNotification)
		userRoutes.PUT("/:username/notifications/notificationId/read", users.AuthMiddleware(server.token), server.UpdateNotification)
		userRoutes.PUT("/:username/notifications/settings", users.AuthMiddleware(server.token), server.UpdateNotificationSettings)
		userRoutes.GET("/:username/notifications/settings", users.AuthMiddleware(server.token), server.GetNotificationSettings)

		// USER SEARCH ENDPOINTS
		userRoutes.GET("/search", users.AuthMiddleware(server.token), server.SearchUsers)

		// USER PRIVACY SETTINGS ENDPOINTS
		userRoutes.GET("/:username/privacy", users.AuthMiddleware(server.token), server.GetPrivacySettings)
		userRoutes.PUT("/:username/privacy", users.AuthMiddleware(server.token), server.UpdatePrivacySettings)

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server.router = r
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
