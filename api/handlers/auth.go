// USER ACCOUNT & AUTHENTICATION ENDPOINTS
package handlers

import (
	"time"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func (server *Server) Register(c *gin.Context) {
	var req users.CreateUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
			"msg":   "Invalid request body",
		})
	}

	if status := users.ValidateRegisterUserReq(req); !status.Valid {
		c.JSON(400, gin.H{
			"msg": status.Message,
		})
	}

	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
			"msg":   "Failed to hash password",
		})
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
	txResult, err := users.CreateUserTx(c, server.db, arg)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
			"msg":   "Failed to create user",
		})
	}

	user := users.UserSerializer{C: c, User: txResult.User}
	c.JSON(200, gin.H{
		"user": user.Response(),
		"msg":  "User created successfully",
	})
}

func (server *Server) ConfirmEmail(c *gin.Context) {
	var req users.ConfirmVerifyEmailParams

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
			"msg":   "Invalid request body",
		})
	}

	if err := users.VerifyEmailTx(c, server.db, users.ConfirmVerifyEmailParams{
		UserID:     req.UserID,
		VerEmailID: req.VerEmailID,
		SecretCode: req.SecretCode,
	}); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
			"msg":   "Failed to verify email",
		})
	}

	c.JSON(200, gin.H{
		"msg": "Email confirmed",
	})
}

// Handlers (simplified)
func (server *Server) Login(c *gin.Context) {
	// Handle user login logic
	// This would include checking user credentials, issuing tokens, etc.
	c.JSON(200, gin.H{
		"message": "Logged in successfully",
	})
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
