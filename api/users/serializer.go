package users

import (
	"github.com/gin-gonic/gin"
)

type UserSerializer struct {
	C *gin.Context
	User
}

type ProfileSerializer struct {
	C *gin.Context
	Profile
}

func (s *UserSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":       s.ID,
		"username": s.Username,
		"email":    s.Email,
		"profile":  s.Profile,
	}
	return response
}

func (s *ProfileSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":       s.ID,
		"username": s.Username,
		"bio":      s.Bio,
		"image":    s.Image,
	}
	return response
}
