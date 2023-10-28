// middlewares.go
package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker, redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			ctx.Abort()
			return
		}

		accessToken := fields[1]

		key := fmt.Sprintf("invalid_tokens:%s", accessToken)
		if redis.Exists(ctx, key).Val() == 1 {
			err := errors.New("invalid access token")
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			ctx.Abort()
			return
		}

		if payload.Valid() != nil {
			utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error(), "Unauthorized")
			ctx.Abort()
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Set("access_token", accessToken)
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

func RateLimitingMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	rate := 10 // Allow 10 requests per second per client IP address.

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Get the current request count for the client IP from Redis
		ctx := context.Background()
		currentCount, err := redisClient.Get(ctx, clientIP).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		// Calculate the remaining requests until the rate limit
		remainingRequests := rate - currentCount

		// Check if the request rate limit has been exceeded
		if remainingRequests < 0 {
			// Set the rate limiting headers
			c.Header("X-RateLimit-Limit", strconv.Itoa(rate))
			c.Header("X-RateLimit-Remaining", "0") // Rate limit exceeded, no remaining requests
			c.Header("X-RateLimit-Reset", "0")     // Reset value is 0 as the limit is exceeded

			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}

		// Set the rate limiting headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remainingRequests))
		c.Header("X-RateLimit-Reset", strconv.Itoa(1)) // Reset value: 1 second

		// Increment the request count in Redis
		_, err = redisClient.Incr(ctx, clientIP).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		// Set an expiration time for the key in Redis (e.g., 1 second)
		if err := redisClient.Expire(ctx, clientIP, time.Second).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

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
