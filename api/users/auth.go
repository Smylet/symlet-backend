package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func (s *UserSerializer) Setup2FA(ctx *gin.Context, db *gorm.DB) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	authPayload, err := GetAuthPayloadFromCtx(ctx)
	if err != nil {
		s.StatusCode = http.StatusUnauthorized
		logger.Error(err.Error())
		return err
	}

	// Find the user based on the authPayload (assuming it contains an ID or other unique identifier).
	var user User
	if err := db.First(&user, "id = ?", authPayload.UserID).Error; err != nil {
		s.StatusCode = http.StatusNotFound
		logger.Error(fmt.Sprintf("User with ID %d not found: %v", authPayload.UserID, err))
		return err
	}

	user.IsTwoFAEnabled = true
	if err := db.Save(&user).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		logger.Error(fmt.Sprintf("Failed to update user 2FA settings: %v", err))
		return err
	}

	return nil

}

func (s *UserSerializer) Verify2FA(ctx *gin.Context, db *gorm.DB, redis *redis.Client, token token.Maker, config utils.Config) error {
	logger := common.NewLogger()

	isVerified, err := verify2FACode(s.TwoFACode, redis, db)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to verify 2FA code")
	}

	if !isVerified {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid 2FA code")
	}

	// Get a list of sessions and get the last one
	var sessions []Session
	if err := db.Model(&Session{}).Where("user_id = ?", s.UserID).Order("created_at desc").Limit(1).Find(&sessions).Error; err != nil {
		return err
	}

	if len(sessions) == 0 {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("session not found")
	}

	// Check the last login IP address - make sure it's the same - if not send an email for attempted login, maybe to change password?
	userIp := ctx.ClientIP()

	if sessions[0].ClientIP != userIp {
		logger.Info("IP address mismatch")

		// Send email for attempted login
		// payload := worker.PayloadSendAttemptedLoginEmail{
		// 	UserName:  sessions[0].UserName,
		// 	Email:     sessions[0].UserName,
		// 	IP:        userIp,
		// 	UserAgent: ctx.GetHeader("User-Agent"),
		// }

		// opts := []asynq.Option{
		// 	asynq.MaxRetry(10),
		// 	asynq.ProcessIn(10 * time.Second),
		// 	asynq.Queue(worker.QueueCritical),
		// }

		// go task.DistributeTaskSendAttemptedLoginEmail(ctx, &payload, opts...)

	}

	accessToken, accessPayload, err := token.CreateToken(s.UserName, s.UserID, config.AccessTokenDuration)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to create access token")
	}

	refreshToken, refreshPayload, err := token.CreateToken(s.UserName, s.UserID, config.RefreshTokenDuration)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to create refresh token")
	}

	s.SessionID = uuid.New()
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.AccessTokenExpiresAt = accessPayload.ExpiresAt
	s.RefreshTokenExpiresAt = refreshPayload.ExpiresAt

	_, err = s.CreateSession(ctx, db, CreateSessionParams{
		ID:           s.SessionID,
		UserID:       s.UserID,
		UserName:     s.UserName,
		RefreshToken: s.RefreshToken,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIP:     ctx.ClientIP(),
		ExpiresAt:    s.AccessTokenExpiresAt,
	})

	if err != nil {
		return err
	}

	return nil
}

var ExpiryTime = time.Hour * 24

func (s *UserSerializer) CreateVerificationEmail(ctx *gin.Context, db *gorm.DB) error {

	verifyEmail := VerificationEmail{
		Email:      s.Email,
		SecretCode: utils.RandomCode(6),
		ExpiresAt:  time.Now().Add(ExpiryTime),
		UserID:     s.UserID,
	}

	if err := db.Create(&verifyEmail).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			s.StatusCode = http.StatusConflict
			return fmt.Errorf("verification email already exists")
		}
		return err
	}

	return nil
}

func (s *UserSerializer) VerifyEmail(ctx *gin.Context, db *gorm.DB) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	err := common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		var verifyEmail VerificationEmail

		if err := tx.Model(&VerificationEmail{}).
			Where("id = ? AND secret_code = ? AND user_id = ?", s.VerEmailID, s.SecretCode, s.UserID).First(&verifyEmail).Error; err != nil {
			if err.Error() == "record not found" {
				s.StatusCode = http.StatusNotFound

				var user User
				if err := tx.Model(&User{}).Where("id = ?", s.UserID).First(&user).Error; err != nil {
					return err
				}

				if user.IsEmailConfirmed {
					return fmt.Errorf("email already confirmed")
				}

				return fmt.Errorf("verification code not found")
			}
			return err
		}

		if verifyEmail.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("verification code has expired")
		}

		if err := tx.Model(&User{}).Where("id = ?", s.UserID).
			Update("is_email_confirmed", true).Error; err != nil {
			return err
		}

		if err := tx.Model(&VerificationEmail{}).Where("id = ?", verifyEmail.ID).Delete(&verifyEmail).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("transaction error: %v", err)
	}

	return nil
}

func (s *UserSerializer) LoginUser(ctx *gin.Context, db *gorm.DB, token token.Maker, config utils.Config, redis *redis.Client) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	err := s.FindByEmail(ctx, db, s.Email)
	if err != nil {
		logger.Error(err.Error())
		if err.Error() == "record not found" {
			s.StatusCode = http.StatusNotFound
			return fmt.Errorf("user not found")
		}
		return err
	}

	if !s.User.IsEmailConfirmed {
		s.StatusCode = http.StatusForbidden
		return fmt.Errorf("email not confirmed")
	}

	if err := common.CheckPassword(s.Password, s.User.Password); err != nil {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid password")
	}

	if s.User.IsTwoFAEnabled {
		TwoFACode, err := generate2FACode(s.User, redis, db)
		if err != nil {
			s.StatusCode = http.StatusInternalServerError
			return fmt.Errorf("failed to login 2FA code")
		}

		s.TwoFACode = TwoFACode

		// payload := worker.PayloadSend2FACode{
		// 	UserName:            s.UserName,
		// 	UserID:              s.UserID,
		// 	VerificationEmailID: s.VerEmailID,
		// 	SecretCode:          s.SecretCode,
		// 	Email:               s.Email,
		// }

		// opts := []asynq.Option{
		// 	asynq.MaxRetry(10),
		// 	asynq.ProcessIn(10 * time.Second),
		// 	asynq.Queue(worker.QueueCritical),
		// }

		// return task.DistributeTaskSend2FACode(ctx, &payload, opts...)

	}

	accessToken, accessPayload, err := token.CreateToken(s.UserName, s.UserID, config.AccessTokenDuration)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to create access token")
	}

	refreshToken, refreshPayload, err := token.CreateToken(s.UserName, s.UserID, config.RefreshTokenDuration)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to create refresh token")
	}

	s.SessionID = uuid.New()
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.AccessTokenExpiresAt = accessPayload.ExpiresAt
	s.RefreshTokenExpiresAt = refreshPayload.ExpiresAt

	_, err = s.CreateSession(ctx, db, CreateSessionParams{
		ID:           s.SessionID,
		UserID:       s.UserID,
		UserName:     s.UserName,
		RefreshToken: s.RefreshToken,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIP:     ctx.ClientIP(),
		ExpiresAt:    s.AccessTokenExpiresAt,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserSerializer) CreateSession(ctx context.Context, db *gorm.DB, arg CreateSessionParams) (Session, error) {
	session := Session{
		ID:           uuid.New(),
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserName:     arg.UserName,
		UserAgent:    arg.UserAgent,
		ClientIP:     arg.ClientIP,
		ExpiresAt:    time.Now().Add(ExpiryTime),
	}

	if err := db.Create(&session).Error; err != nil {
		return Session{}, err
	}

	return session, nil

}

func (s *UserSerializer) FindSession(ctx context.Context, db *gorm.DB, UserID uint) (Session, error) {
	var session Session

	if err := db.Model(&Session{}).Where("user_id = ?", UserID).First(&session).Error; err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s *UserSerializer) RenewAccessToken(ctx *gin.Context, db *gorm.DB, token token.Maker, config utils.Config) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	refreshPayload, err := token.VerifyToken(s.RefreshToken)
	if err != nil {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid refresh token")
	}

	if refreshPayload.ExpiresAt.Before(time.Now()) {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("refresh token has expired")
	}

	session, err := s.FindSession(ctx, db, refreshPayload.UserID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if session.RefreshToken != s.RefreshToken {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid refresh token")
	}

	if session.IsBlocked {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("session is blocked")
	}

	if session.UserName != refreshPayload.Username {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("incorrect session user")
	}

	accessToken, accessPayload, err := token.CreateToken(session.UserName, session.UserID, config.AccessTokenDuration)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to renew access token")
	}

	s.AccessToken = accessToken
	s.AccessTokenExpiresAt = accessPayload.ExpiresAt

	return nil
}

func (s *UserSerializer) ForgotPassword(ctx *gin.Context, db *gorm.DB, task worker.TaskDistributor, token token.Maker) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	err := s.FindByEmail(ctx, db, s.Email)
	if err != nil {
		logger.Error(err.Error())
		if err.Error() == "record not found" {
			s.StatusCode = http.StatusNotFound
			return fmt.Errorf("user not found")
		}
		return err
	}

	if !s.User.IsEmailConfirmed {
		s.StatusCode = http.StatusForbidden
		return fmt.Errorf("email not confirmed")
	}

	ExpiryTime := time.Hour * 24

	resetToken, resetTokenPayload, err := token.CreateToken(s.UserName, s.UserID, ExpiryTime)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("failed to create access token")
	}

	s.ResetPasswordToken = resetToken
	s.ResetPasswordTokenExpiresAt = resetTokenPayload.ExpiresAt

	updates := map[string]interface{}{
		"reset_password_token":      resetToken,
		"reset_password_expires_at": resetTokenPayload.ExpiresAt,
	}

	// Update the specified fields for the user with the given ID.
	if err := db.Model(&User{}).Where("id = ?", s.UserID).Updates(updates).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		return err
	}

	payload := worker.PayloadSendForgetPasswordEmail{
		UserName:   s.UserName,
		Email:      s.Email,
		ResetToken: s.ResetPasswordToken,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	go func() {
		err := task.DistributeTaskSendForgetPasswordEmail(ctx, &payload, opts...)
		if err != nil {
			// Handle the error here, e.g., log it
			logger.Error(err.Error())
		}
	}()

	// User clicks the link and is redirected to a "Reset Password" page - Client GET /reset-password?token=unique-reset-token
	// The "Reset Password" page render a password reset form for the user to enter a new password.
	return nil
}

func (s *UserSerializer) PasswordReset(ctx *gin.Context, db *gorm.DB, token token.Maker) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	resetPayload, err := token.VerifyToken(s.ResetPasswordToken)
	if err != nil {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid reset password token")
	}

	if resetPayload.ExpiresAt.Before(time.Now()) {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("reset password token has expired")
	}

	var user User
	if err := db.Model(&User{}).Where("reset_password_token = ?", s.ResetPasswordToken).First(&user).Error; err != nil {
		logger.Error(err.Error())
		if err.Error() == "record not found" {
			s.StatusCode = http.StatusBadRequest
			return fmt.Errorf("invalid reset password token")
		}
		return err
	}

	hashedPassword, err := common.HashPassword(s.Password)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		logger.Error(err.Error())
		return err
	}

	user.Password = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetPasswordExpiresAt = time.Time{}

	if err := db.Save(&user).Error; err != nil {
		logger.Error(err.Error())
		return err
	}

	s.UserID = user.ID
	s.Email = user.Email
	s.UserName = user.UserName
	s.User = &user
	return nil
}

func (s *UserSerializer) ChangePassword(ctx *gin.Context, db *gorm.DB, task worker.TaskDistributor) error {

	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	authPayload, err := GetAuthPayloadFromCtx(ctx)
	if err != nil {
		s.StatusCode = http.StatusUnauthorized
		logger.Error(err.Error())
		return err
	}

	var user User
	if err := db.Model(&User{}).Where("id = ?", authPayload.UserID).First(&user).Error; err != nil {
		s.StatusCode = http.StatusUnauthorized
		logger.Error(err.Error())
		return err
	}

	if err := common.CheckPassword(s.Password, user.Password); err != nil {
		s.StatusCode = http.StatusUnauthorized
		return fmt.Errorf("invalid password")
	}

	hashedPassword, err := common.HashPassword(s.NewPassword)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		logger.Error(err.Error())
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		logger.Error(err.Error())
		return err
	}

	payload := worker.PayloadSendChangePasswordEmail{
		UserName: user.UserName,
		Email:    user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueDefault),
	}

	go func() {
		err := task.DistributeTaskSendChangePasswordEmail(ctx, &payload, opts...)
		if err != nil {
			// Handle the error here, e.g., log it
			logger.Error(err.Error())
		}
	}()

	s.UserID = user.ID
	s.Email = user.Email
	s.UserName = user.UserName
	s.User = &user

	return nil
}

func (s *UserSerializer) ResendEmailVerification(ctx *gin.Context, db *gorm.DB, task worker.TaskDistributor) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	err := s.FindByEmail(ctx, db, s.Email)
	if err != nil {
		logger.Error(err.Error())
		if err.Error() == "record not found" {
			s.StatusCode = http.StatusNotFound
			return fmt.Errorf("user not found")
		}
		return err
	}

	if s.User.IsEmailConfirmed {
		s.StatusCode = http.StatusForbidden
		return fmt.Errorf("email already confirmed")
	}

	err = s.CreateVerificationEmail(ctx, db)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	payload := worker.PayloadSendVerifyEmail{
		UserName:            s.UserName,
		UserID:              s.UserID,
		VerificationEmailID: s.VerEmailID,
		SecretCode:          s.SecretCode,
		Email:               s.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	return task.DistributeTaskSendVerifyEmail(ctx, &payload, opts...)

}

func (s *UserSerializer) LogoutUser(ctx *gin.Context, db *gorm.DB, redis *redis.Client) error {

	accessToken, exist := ctx.Get("access_token")
	if !exist {
		s.StatusCode = http.StatusBadRequest
		return fmt.Errorf("access token does not exist")
	}

	key := fmt.Sprintf("invalid_tokens:%s", accessToken)
	if err := redis.Set(context.Background(), key, "logged out", time.Hour*24).Err(); err != nil {
		s.StatusCode = http.StatusInternalServerError
		return err
	}
	return nil
}
