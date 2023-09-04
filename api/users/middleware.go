// middlewares.go
package users

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

// 2. Authorization Middleware
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add user role to the payload
		// check the role in this middleware
		// Implement your authorization logic here. E.g., check the user role or specific permissions.
		if /* user does not have permission */ false {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 3. Rate Limiting Middleware
func RateLimitingMiddleware() gin.HandlerFunc {
	// This is a very basic implementation. You might want to use a more sophisticated approach with Redis or another store.
	rate := time.Second / 10 // For demonstration: Limit to 10 requests per second.
	lastRequest := time.Now()

	return func(c *gin.Context) {
		if time.Since(lastRequest) < rate {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		lastRequest = time.Now()
		c.Next()
	}
}

// 4. Error Handling Middleware
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				// You can also log the error if needed.
			}
		}()
		c.Next()
	}
}

// 6. Two-Factor Authentication Middleware
func TwoFactorAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has 2FA enabled.
		// If so, verify the provided 2FA code.
		if /* 2FA is enabled and the code is incorrect */ false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2FA verification failed"})
			c.Abort()
			return
		}
		c.Next()
	}
}
