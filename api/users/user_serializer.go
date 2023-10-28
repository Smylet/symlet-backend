package users

import (
	"context"
	"time"

	// "github.com/Smylet/symlet-backend/utilities/worker"

	"github.com/google/uuid"
)

type UserRepositoryProvider interface {
	Create(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	Get(ctx context.Context, arg FindUserParams) (UserSerializer, error)
	Update(ctx context.Context, arg UpdateUserParams) (UserSerializer, error)
	CreateVerificationEmail(ctx context.Context, req CreateVerificationEmailParams) (VerificationEmail, error)
	VerifyEmail(ctx context.Context, arg ConfirmVerificationEmailParams) error
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
}

type UserSerializer struct {
	UserName                    string `json:"user_name" form:"user_name" custom_binding:"requiredFor:create" validate:"min=3"`
	Email                       string `json:"email" form:"email" custom_binding:"requiredFor:create,update,forgotPassword,resendEmailConfirmation" validate:"email"`
	Password                    string `json:"password" form:"password" custom_binding:"requiredFor:create,resetPassword,changePassword" validate:"min=8"`
	SecretCode                  uint   `json:"secret_code" form:"secret_code" custom_binding:"requiredFor:verifyEmail"`
	VerEmailID                  uint   `json:"ver_email_id" form:"ver_email_id" custom_binding:"requiredFor:verifyEmail"`
	UserID                      uint   `json:"user_id" form:"user_id" custom_binding:"requiredFor:verifyEmail"`
	AccessToken                 string `json:"access_token"`
	RefreshToken                string `json:"refresh_token" form:"refresh_token" custom_binding:"requiredFor:renewAccessToken"`
	SessionID                   uuid.UUID
	AccessTokenExpiresAt        time.Time              `json:"access_token_expires_at"`
	RefreshTokenExpiresAt       time.Time              `json:"refresh_token_expires_at"`
	User                        *User                  `json:"-"`
	StatusCode                  int                    `json:"-"`
	Page                        int                    `json:"page" form:"page" custom_binding:"requiredFor:list"`
	Limit                       int                    `json:"limit" form:"limit" custom_binding:"requiredFor:list"`
	ResetPasswordToken          string                 `json:"reset_password_token" form:"reset_password_token" custom_binding:"requiredFor:resetPassword"`
	ResetPasswordTokenExpiresAt time.Time              `json:"reset_password_token_expires_at"`
	NewPassword                 string                 `json:"new_password" form:"new_password" custom_binding:"requiredFor:changePassword"`
	TwoFACode                   uint                   `json:"two_fa_code" form:"two_fa_code" custom_binding:"requiredFor:verify2FA"`
	Preferences                 map[string]interface{} `json:"preferences" form:"preferences" custom_binding:"requiredFor:updatePreferences"`
	PastSearches                []PastSearch           `json:"past_searches" form:"past_searches" custom_binding:"requiredFor:addPastSearches"`
}
