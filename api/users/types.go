// validators.go
package users

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginReq struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID      `json:"session_id"`
	AccessToken           string         `json:"access_token"`
	AccessTokenExpiresAt  time.Time      `json:"access_token_expires_at"`
	RefreshToken          string         `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time      `json:"refresh_token_expires_at"`
	User                  UserSerializer `json:"user"`
}
type CreateUserTxParams struct {
	CreateUserReq
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

type CreateVerifyEmailParams struct {
	UserID     uint
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

type ConfirmVerifyEmailParams struct {
	UserID     uint   `form:"user_id" binding:"required"`
	VerEmailID uint   `form:"ver_email_id" binding:"required"`
	SecretCode string `form:"secret_code" binding:"required"`
}

type ValidationStatus struct {
	Valid   bool
	Message string
}
type UpdateVerifyEmailParams struct {
	Email      string
	SecretCode string
}

type CreateSessionParams struct {
	ID           uuid.UUID
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiresAt    time.Time
}
