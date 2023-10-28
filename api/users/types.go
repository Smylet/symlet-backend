// validators.go
package users

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
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
type FindUserParams struct {
	User
	IncludeProfile bool
}
type UpdateUserParams struct {
	Criteria User // Fields to search for the user
	Updates  User // Fields to update for the user
}

type CreateUserTxParams struct {
	CreateUserReq
	AfterCreate func() error
}

type CreateUserTxResult struct {
	User User
}

type ConfirmVerificationEmailParams struct {
	UserID     uint
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

type CreateVerificationEmailParams struct {
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
	UserID       uint
	UserName     string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiresAt    time.Time
}
