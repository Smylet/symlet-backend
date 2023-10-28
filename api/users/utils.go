package users

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/Smylet/symlet-backend/utilities/token"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func buildSearchQuery(database *gorm.DB, filter User) *gorm.DB {
	query := database.Model(&User{})

	// Reflect on the User type to iterate over its fields
	val := reflect.ValueOf(filter)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() { // if the field is set
			// Note: This assumes the struct field name and the DB column name are the same.
			// You might need to handle this differently in production code.
			dbFieldName := val.Type().Field(i).Name
			query = query.Where(fmt.Sprintf("%s = ?", dbFieldName), field.Interface())
		}
	}
	return query
}

func GetAuthPayloadFromCtx(ctx *gin.Context) (*token.Payload, error) {
	payload, exist := ctx.Get(AuthorizationPayloadKey)
	if !exist {
		return nil, fmt.Errorf("authorization payload does not exist")
	}
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("authorization payload is not of type *token.Payload")
	}
	return authPayload, nil
}

func generate2FACode(user *User, redis *redis.Client, db *gorm.DB) (uint, error) {
	// Generate a random six-digit code using crypto/rand
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	diff := new(big.Int).Sub(max, min)
	codeBig, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, errors.New("failed to generate secure random number")
	}
	code := uint(codeBig.Add(codeBig, min).Uint64())

	// Set the key as a combination of a constant prefix and the user's ID
	// This makes it easy to fetch the code by user
	key := fmt.Sprintf("user:%d:2fa", user.ID)

	// store the code in db
	user.TwoFASecret = code
	err = db.Save(&user).Error
	if err != nil {
		return 0, errors.New("failed to set 2FA code in db")
	}

	// Store the code in Redis with a 5-minute expiration
	err = redis.Set(context.TODO(), key, code, 5*time.Minute).Err()
	if err != nil {
		return 0, errors.New("failed to set 2FA code in Redis")
	}

	return code, nil
}

func verify2FACode(enteredCode uint, redisClient *redis.Client, db *gorm.DB) (bool, error) {
	var user User

	// Get the user with the entered code
	if err := db.First(&user, "two_fa_secret = ?", enteredCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("2FA code is invalid")
		}
		return false, fmt.Errorf("failed to get user with 2FA code: %v", err)
	}

	key := fmt.Sprintf("user:%d:2fa", user.ID)

	// Get the stored code from Redis
	storedCode, err := redisClient.Get(context.TODO(), key).Uint64()
	if err != nil {
		// Handle case where the code is not found (expired or never set)
		if errors.Is(err, redis.Nil) {
			return false, errors.New("2FA code has expired or does not exist")
		}
		return false, fmt.Errorf("failed to get 2FA code from Redis: %v", err)
	}

	// Check if the entered code matches the stored code
	if uint(storedCode) == enteredCode {
		// If the code is correct, you might want to delete it from Redis to ensure it can't be used again.
		// But this is optional and depends on your requirements.
		_ = redisClient.Del(context.TODO(), key)
		return true, nil
	}

	return false, nil
}
