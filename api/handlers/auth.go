// USER ACCOUNT & AUTHENTICATION ENDPOINTS
package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (server *Server) Register(c *gin.Context) {
	log.Println("here")

	// user := users.CreateUserTx()

	// Handle user registration logic
	// This would include creating a new user account, sending a confirmation email, etc.
	c.JSON(200, gin.H{
		"message": "Registered successfully",
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

func (server *Server) ConfirmEmail(c *gin.Context) {
	// Handle email confirmation logic
	c.JSON(200, gin.H{
		"message": "Email confirmed",
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
