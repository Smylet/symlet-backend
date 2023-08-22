// validators.go
package users

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserRegistrationValidator struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (v *UserRegistrationValidator) Bind(c *gin.Context) error {
	// ... other code ...

	// Check if password contains at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(v.Password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check if password contains at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(v.Password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check if password contains at least one digit
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(v.Password) {
		return errors.New("password must contain at least one digit")
	}

	// Check if password contains at least one special character
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)
	if !specialCharRegex.MatchString(v.Password) {
		return errors.New("password must contain at least one special character")
	}

	// Check if password is at least 8 characters long
	if len(v.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
