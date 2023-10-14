package handlers

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/users"
	_ "github.com/Smylet/symlet-backend/docs"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/sms"
	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
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
	sms *sms.SMSSender
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config, db *gorm.DB, task worker.TaskDistributor, mailer mail.EmailSender, session *session.Session, sms *sms.SMSSender) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config: config,
		cron:   cron.New(),
		db:     db,
		task:   task,
		mailer: mailer,
		token:  tokenMaker,
		sms: sms,
		
	}
	if config.Environment == "production" {
		server.session = session
	}

	gin.SetMode(gin.ReleaseMode)

	server.registerRoutes()

	return server, nil
}

func (server *Server) registerRoutes() {
	r := gin.Default()
	hostelRoutes := r.Group("/hostels")
	{
		hostelRoutes.POST("/", users.AuthMiddleware(server.token), server.CreateHostel)
		hostelRoutes.GET("/:uid", server.GetHostel)
		hostelRoutes.GET("/", server.ListHostels)
		hostelRoutes.PATCH("/:uid", users.AuthMiddleware(server.token), server.UpdateHostel)
		hostelRoutes.DELETE("/:uid", users.AuthMiddleware(server.token), server.DeleteHostel)
	}

	managerRoutes := r.Group("/hostel-managers")
	{
		managerRoutes.POST("/", users.AuthMiddleware(server.token), server.CreateHostelManager)
	}

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
		//userRoutes.POST("/profile", users.AuthMiddleware(server.token), server.CreateUserProfile)
		userRoutes.GET("/profile/:uid", users.AuthMiddleware(server.token), server.GetUserProfile)
		userRoutes.PUT("/profile/:uid", users.AuthMiddleware(server.token), server.EditUserProfile)
		userRoutes.DELETE("/:uid", users.AuthMiddleware(server.token), server.DeleteUserProfile)
		userRoutes.GET("/:uid/history", users.AuthMiddleware(server.token), server.ViewProfileEditHistory)
		userRoutes.GET("/:uid/backups", users.AuthMiddleware(server.token), server.ListProfileBackups)
		userRoutes.PUT("/:uid/restore", users.AuthMiddleware(server.token), server.RestoreUserProfile)
		userRoutes.GET("/:uid/export", users.AuthMiddleware(server.token), server.ExportUserProfile)
		userRoutes.PUT("/:uid/deactivate", users.AuthMiddleware(server.token), server.DeactivateAccount)
		userRoutes.PUT("/:uid/reactivate", users.AuthMiddleware(server.token), server.ReactivateAccount)

		// // USER PROFILE PICTURE ENDPOINTS
		userRoutes.POST("/:uid/picture", users.AuthMiddleware(server.token), server.UploadProfilePicture)
		userRoutes.PUT("/:uid/picture", users.AuthMiddleware(server.token), server.UpdateProfilePicture)
		userRoutes.DELETE("/:uid/picture", users.AuthMiddleware(server.token), server.DeleteProfilePicture)

		//USER NOTIFICATION ENDPOINTS
		userRoutes.GET("/:uid/notifications", users.AuthMiddleware(server.token), server.GetNotification)
		userRoutes.PUT("/:uid/notifications/notificationId/read", users.AuthMiddleware(server.token), server.UpdateNotification)
		userRoutes.PUT("/:uid/notifications/settings", users.AuthMiddleware(server.token), server.UpdateNotificationSettings)
		userRoutes.GET("/:uid/notifications/settings", users.AuthMiddleware(server.token), server.GetNotificationSettings)

		// USER SEARCH ENDPOINTS
		userRoutes.GET("/search", users.AuthMiddleware(server.token), server.SearchUsers)

		// // USER PRIVACY SETTINGS ENDPOINTS
		userRoutes.GET("/:uid/privacy", users.AuthMiddleware(server.token), server.GetPrivacySettings)
		userRoutes.PUT("/:uid/privacy", users.AuthMiddleware(server.token), server.UpdatePrivacySettings)

	}
	referenceRoutes := r.Group("/references")
	{
		referenceRoutes.GET("/amenities", server.ListAmenities)
		referenceRoutes.GET("/universities", server.ListUniversities)
	}
	studentRoutes := r.Group("/students")
	{
		studentRoutes.POST("/", users.AuthMiddleware(server.token), server.CreateStudent)
	}
	vendorRoutes := r.Group("/vendors")
	{
		vendorRoutes.POST("/", users.AuthMiddleware(server.token), server.CreateVendor)
		vendorRoutes.GET("/:uid", server.GetVendor)
	}
	userRoute := r.Group("/user")
	{
		// USER 2FA ENDPOINTS
		userRoute.POST("/2fa/setup", server.Setup2FA)
		userRoute.POST("/2fa/verify", server.Verify2FA)
	}

	// health check endpoint
	r.GET("/health", server.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server.router = r
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	address = strings.TrimPrefix(address, "http://")
	return server.router.Run(address)
}
