// models.go
package users

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Smylet/symlet-backend/utilities/common"
	"gorm.io/gorm"
)

type User struct {
	common.AbstractBaseModel
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `sql:"not null"`
	Profile   Profile
}

type Profile struct {
	ID     uint `gorm:"primary_key"`
	UserID uint
	Bio    string
	Image  string
}

type VerifyEmail struct {
	ID         uint `gorm:"primary_key"`
	Email      string
	SecretCode string
	ExpiresAt  time.Time
}

func CreateUserTx(ctx context.Context, database *gorm.DB, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := common.ExecTx(ctx, database, func(tx *gorm.DB) error {
		user := User{
			Username:  arg.Username,
			Email:     arg.Email,
			Password:  arg.Password,
			CreatedAt: time.Now(),
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

func FindUser(ctx context.Context, database *gorm.DB, arg FindUserParams) (UserSerializer, error) {
	var result UserSerializer

	query := database.WithContext(ctx).Model(&User{})

	v := reflect.ValueOf(arg.User)
	t := reflect.TypeOf(arg.User)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Only process non-zero fields.
		if !field.IsZero() {
			// Assumes that the column name in the DB is same as struct field name.
			// If they're different, you'd use gorm struct tags to get DB column names.
			query = query.Where(fmt.Sprintf("%s = ?", fieldType.Name), field.Interface())
		}
	}

	if arg.IncludeProfile {
		query = query.Preload("Profile")
	}

	if err := query.First(&result.User).Error; err != nil {
		return UserSerializer{}, fmt.Errorf("failed to get user: %w", err)
	}

	return result, nil
}

func CreateVerifyEmail(ctx context.Context, database *gorm.DB, req CreateVerifyEmailParams) (VerifyEmail, error) {
	verifyEmail := VerifyEmail{
		Email:      req.Email,
		SecretCode: req.SecretCode,
		ExpiresAt:  time.Now().Add(time.Hour * 24),
	}

	if err := database.Create(&verifyEmail).Error; err != nil {
		return verifyEmail, err
	}

	return verifyEmail, nil
}

func UpdateVerifyEmail(ctx context.Context, database *gorm.DB, req UpdateVerifyEmailParams) (VerifyEmail, error) {
	var verifyEmail VerifyEmail
	if err := database.Where("email = ?", req.Email).First(&verifyEmail).Error; err != nil {
		return verifyEmail, err
	}

	verifyEmail.SecretCode = req.SecretCode

	if err := database.Save(&verifyEmail).Error; err != nil {
		return verifyEmail, err
	}

	return verifyEmail, nil
}

// type VerifyEmailTxParams struct {
// 	EmailId    int64
// 	SecretCode string
// }

// type VerifyEmailTxResult struct {
// 	User User
// 	err  error // Add an error field of type error
// }

// func VerifyEmailTx(ctx context.Context, database *gorm.DB, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
// 	var result VerifyEmailTxResult

// 	err := common.ExecTx(ctx, database, func(tx *gorm.DB) error {

// 		return nil
// 	})
// 	if err != nil {
// 		return VerifyEmailTxResult{}, fmt.Errorf("transaction error: %v", err)
// 	}

// 	return result, nil
// }
