package users

import (
	"github.com/gin-gonic/gin"
)

type UserSerializer struct {
	C *gin.Context
	User
	Profile ProfileSerializer `json:"profile"`
}

type ProfileSerializer struct {
	C *gin.Context
	Profile
}

func (s *UserSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":         s.ID,
		"username":   s.Username,
		"email":      s.Email,
		"created_at": s.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"profile":    s.Profile.Response(),
	}
	return response
}

func (s *ProfileSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":    s.ID,
		"bio":   s.Bio,
		"image": s.Image,
	}
	return response
}

type VerifyEmailRequest struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // This can be used to send the expiration time of the token, if needed.
}
