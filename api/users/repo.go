package users

import (
	"context"
	"fmt"
	"time"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryProvider interface {
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	FindUser(ctx context.Context, arg FindUserParams) (UserSerializer, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (UserSerializer, error)
	CreateVerifyEmail(ctx context.Context, req CreateVerifyEmailParams) (VerificationEmail, error)
	VerifyEmailTx(ctx context.Context, arg ConfirmVerifyEmailParams) error
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
}

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepositoryProvider {
	return &UserRepository{database}
}

func (user UserRepository) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := common.ExecTx(ctx, user.database, func(tx *gorm.DB) error {
		user := User{
			Username: arg.Username,
			Email:    arg.Email,
			Password: arg.Password,
		}
		if err := tx.Create(&user).Error; err != nil {
			// check if the error is a duplicate key error
			return err
		}

		profile := Profile{
			UserID: user.ID,
		}
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}
		result.User = user

		return arg.AfterCreate(result.User)
	})
	if err != nil {
		return CreateUserTxResult{}, fmt.Errorf("transaction error: %v", err)
	}

	return result, nil
}

type FindUserParams struct {
	User
	IncludeProfile bool
}

func (user UserRepository) FindUser(ctx context.Context, arg FindUserParams) (UserSerializer, error) {
	var result UserSerializer

	query := buildQueryFromUser(ctx, user.database, Args{User: arg.User})

	// if arg.IncludeProfile {
	// 	query = query.Preload("Profile")
	// }

	if err := query.First(&result.User).Error; err != nil {
		return UserSerializer{}, fmt.Errorf("failed to get user: %w", err)
	}

	return result, nil
}

type UpdateUserParams struct {
	Criteria User // Fields to search for the user
	Updates  User // Fields to update for the user
}

func (user UserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) (UserSerializer, error) {
	var result UserSerializer

	// Use the buildQueryFromUser function to get a query based on non-zero fields.
	query := buildQueryFromUser(ctx, user.database, Args{User: arg.Criteria})

	arg.Updates.UpdatedAt = time.Now()
	query = query.Updates(arg.Updates)

	// Execute the updates based on the built query.

	if err := query.First(&result.User).Error; err != nil {
		return UserSerializer{}, fmt.Errorf("failed to retrieve updated user: %w", err)
	}

	return result, nil
}

var ExpiryTime = time.Hour * 24

func (user UserRepository) CreateVerifyEmail(ctx context.Context, req CreateVerifyEmailParams) (VerificationEmail, error) {
	verifyEmail := VerificationEmail{
		UserID:     req.UserID,
		Email:      req.Email,
		SecretCode: req.SecretCode,
		ExpiresAt:  time.Now().Add(ExpiryTime),
	}

	if err := user.database.Create(&verifyEmail).Error; err != nil {
		return verifyEmail, err
	}

	return verifyEmail, nil
}

func (user UserRepository) VerifyEmailTx(ctx context.Context, arg ConfirmVerifyEmailParams) error {
	return common.ExecTx(ctx, user.database, func(tx *gorm.DB) error {
		var verifyEmail VerificationEmail

		// Chain the query conditions and fetch the record in one operation
		if err := tx.Model(&VerificationEmail{}).
			Where("user_id = ? AND id = ?", arg.UserID, arg.VerEmailID).
			First(&verifyEmail).Error; err != nil {
			return err
		}

		// Check the secret code and expiration
		if verifyEmail.SecretCode != arg.SecretCode {
			return fmt.Errorf("invalid secret code")
		}
		if verifyEmail.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("expired secret code")
		}

		// Update the user's email confirmation status
		if err := tx.Model(&User{}).
			Where("id = ?", arg.UserID).
			Update("is_email_confirmed", true).Error; err != nil {
			return err
		}

		// Delete the verification record
		if err := tx.Delete(&verifyEmail).Error; err != nil {
			return err
		}

		return nil
	})
}

func (user UserRepository) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	session := Session{
		ID:           uuid.New(),
		Username:     arg.Username,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIP:     arg.ClientIP,
		ExpiresAt:    time.Now().Add(ExpiryTime),
	}

	if err := user.database.Create(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}
