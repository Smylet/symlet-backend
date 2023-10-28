package handlers

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/users"
	_ "github.com/Smylet/symlet-backend/docs"
	"github.com/Smylet/symlet-backend/utilities/jobs"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
)

type Server struct {
	router      *gin.Engine
	config      utils.Config
	cron        *cron.Cron
	db          *gorm.DB
	task        worker.TaskDistributor
	mailer      mail.EmailSender
	session     *session.Session
	token       token.Maker
	redisClient *redis.Client
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config, db *gorm.DB, task worker.TaskDistributor, mailer mail.EmailSender, session *session.Session, redis *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:      config,
		cron:        cron.New(),
		db:          db,
		task:        task,
		mailer:      mailer,
		token:       tokenMaker,
		redisClient: redis,
	}
	if config.Environment == "production" {
		server.session = session
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	redisStore := persist.NewRedisStore(redis)

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cache.CacheByRequestURI(redisStore, 2*time.Second))
	r.Use(users.RateLimitingMiddleware(redis))

	server.registerRoutes(r)

	if err = jobs.InitializeCronJobs(server.cron, server.db, server.task); err != nil {
		return nil, err
	}
	return server, nil
}

func (server *Server) registerRoutes(r *gin.Engine) {

	hostelRoutes := r.Group("/hostels")
	{
		hostelRoutes.POST("/", users.AuthMiddleware(server.token, server.redisClient), server.CreateHostel)
		hostelRoutes.GET("/:uid", server.GetHostel)
		hostelRoutes.GET("/", server.ListHostels)
		hostelRoutes.PATCH("/:uid", users.AuthMiddleware(server.token, server.redisClient), server.UpdateHostel)
		hostelRoutes.DELETE("/:uid", users.AuthMiddleware(server.token, server.redisClient), server.DeleteHostel)
	}

	managerRoutes := r.Group("/hostel-managers")
	{
		managerRoutes.POST("/", users.AuthMiddleware(server.token, server.redisClient), server.CreateHostelManager)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", server.CreateUser)
		userRoutes.POST("/login", server.LoginUser)
		userRoutes.GET(":uid", users.AuthMiddleware(server.token, server.redisClient), server.GetUser)
		userRoutes.GET("/me", users.AuthMiddleware(server.token, server.redisClient), server.GetMe)
		userRoutes.POST("/logout", users.AuthMiddleware(server.token, server.redisClient), server.LogoutUser)
		userRoutes.POST("/confirm-email", server.verifyEmail)
		userRoutes.POST("/email/confirmation/resend", server.ResendEmailVerification)
		userRoutes.GET("/", users.AuthMiddleware(server.token, server.redisClient), server.GetUsers)

		userRoutes.GET("/:uid/preferences", users.AuthMiddleware(server.token, server.redisClient), server.GetPreferences)
		userRoutes.PUT("/:uid/preferences", users.AuthMiddleware(server.token, server.redisClient), server.UpdatePreferences)
		userRoutes.DELETE("/:uid/preferences", users.AuthMiddleware(server.token, server.redisClient), server.DeletePreferences)
		userRoutes.GET("/:uid/past-searches", users.AuthMiddleware(server.token, server.redisClient), server.GetPastSearches)
		userRoutes.POST("/:uid/past-searches", users.AuthMiddleware(server.token, server.redisClient), server.AddPastSearch)
		userRoutes.DELETE("/:uid/past-searches", users.AuthMiddleware(server.token, server.redisClient), server.ClearPastSearches)

		// USER SEARCH ENDPOINTS
		userRoutes.GET("/search", users.AuthMiddleware(server.token, server.redisClient), server.SearchUsers)

	}

	profileRoutes := r.Group("/profile")
	{
		// USER PROFILE MANAGEMENT ENDPOINTS
		profileRoutes.GET("/:uid", users.AuthMiddleware(server.token, server.redisClient), server.GetUserProfile)
		profileRoutes.PUT("/:uid", users.AuthMiddleware(server.token, server.redisClient), server.EditUserProfile)
		profileRoutes.DELETE("/:uid", users.AuthMiddleware(server.token, server.redisClient), server.DeleteUserProfile)
		profileRoutes.GET("/:uid/history", users.AuthMiddleware(server.token, server.redisClient), server.ViewProfileEditHistory)
		profileRoutes.GET("/:uid/backups", users.AuthMiddleware(server.token, server.redisClient), server.ListProfileBackups)
		profileRoutes.PUT("/:uid/restore", users.AuthMiddleware(server.token, server.redisClient), server.RestoreUserProfile)
		profileRoutes.GET("/:uid/export", users.AuthMiddleware(server.token, server.redisClient), server.ExportUserProfile)
		profileRoutes.PUT("/:uid/deactivate", users.AuthMiddleware(server.token, server.redisClient), server.DeactivateAccount)
		profileRoutes.PUT("/:uid/reactivate", users.AuthMiddleware(server.token, server.redisClient), server.ReactivateAccount)

		// // USER PROFILE PICTURE ENDPOINTS
		profileRoutes.POST("/:uid/picture", users.AuthMiddleware(server.token, server.redisClient), server.UploadProfilePicture)
		profileRoutes.PUT("/:uid/picture", users.AuthMiddleware(server.token, server.redisClient), server.UpdateProfilePicture)
		profileRoutes.DELETE("/:uid/picture", users.AuthMiddleware(server.token, server.redisClient), server.DeleteProfilePicture)

		// USER NOTIFICATION ENDPOINTS
		profileRoutes.GET("/:uid/notifications", users.AuthMiddleware(server.token, server.redisClient), server.GetNotification)
		profileRoutes.PUT("/:uid/notifications/notificationId/read", users.AuthMiddleware(server.token, server.redisClient), server.UpdateNotification)
		profileRoutes.PUT("/:uid/notifications/settings", users.AuthMiddleware(server.token, server.redisClient), server.UpdateNotificationSettings)
		profileRoutes.GET("/:uid/notifications/settings", users.AuthMiddleware(server.token, server.redisClient), server.GetNotificationSettings)

		// // USER PRIVACY SETTINGS ENDPOINTS
		profileRoutes.GET("/:uid/privacy", users.AuthMiddleware(server.token, server.redisClient), server.GetPrivacySettings)
		profileRoutes.PUT("/:uid/privacy", users.AuthMiddleware(server.token, server.redisClient), server.UpdatePrivacySettings)

	}

	referenceRoutes := r.Group("/references")
	{
		referenceRoutes.GET("/amenities", server.ListAmenities)
		referenceRoutes.GET("/universities", server.ListUniversities)
	}

	studentRoutes := r.Group("/students")
	{
		studentRoutes.POST("/", users.AuthMiddleware(server.token, server.redisClient), server.CreateStudent)
	}

	userRoute := r.Group("/user")
	{
		// USER 2FA ENDPOINTS
		userRoute.POST("/2fa/setup", server.Setup2FA)
		userRoute.POST("/2fa/verify", server.Verify2FA)
	}

	// health check endpoint
	r.GET("/health", server.HealthCheck)

	// AUTHENTICATION ENDPOINTS
	r.POST("/token/renew", server.RenewAccessToken)
	r.POST("/forgot-password", server.ForgotPassword)
	r.POST("/reset-password", server.PasswordReset)
	r.POST("/change-password", users.AuthMiddleware(server.token, server.redisClient), server.ChangePassword)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server.router = r
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	address = strings.TrimPrefix(address, "http://")
	return server.router.Run(address)
}
