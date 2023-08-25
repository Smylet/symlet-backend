package handlers

import "github.com/gin-gonic/gin"

func (server *Server) GetNotification(c *gin.Context) {
	// Get user ID from JWT
	// Get user from DB
	// Return user's notification settings
}

func (server *Server) UpdateNotification(c *gin.Context) {
	// Get user ID from JWT
	// Get user from DB
	// Update user's notification settings
	// Return updated user
}

func (server *Server) UpdateNotificationSettings(c *gin.Context) {
	// Get user ID from JWT
	// Get user from DB
	// Update user's notification settings
	// Return updated user
}

func (server *Server) GetNotificationSettings(c *gin.Context) {
	// Get user ID from JWT
	// Get user from DB
	// Return user's notification settings
}
