package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func (s *UserSerializer) Create(ctx *gin.Context, db *gorm.DB, task worker.TaskDistributor, cron *cron.Cron) error {
	logger := common.NewLogger()

	if err := s.Validate(); err != nil {
		s.SecretCode = http.StatusBadRequest
		return fmt.Errorf("validation error: %v", err)
	}

	hashedPassword, err := common.HashPassword(s.Password)
	if err != nil {
		s.StatusCode = http.StatusInternalServerError
		logger.Error(err.Error())
		return err
	}

	arg := CreateUserTxParams{
		CreateUserReq: CreateUserReq{
			UserName: s.UserName,
			Email:    s.Email,
			Password: hashedPassword,
		},
		AfterCreate: func() error {
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
		},
	}

	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		user := User{
			UserName: arg.UserName,
			Email:    arg.Email,
			Password: arg.Password,
		}

		if err := tx.Create(&user).Error; err != nil {
			s.StatusCode = http.StatusInternalServerError
			if strings.Contains(err.Error(), "duplicate") {
				s.StatusCode = http.StatusConflict
				return fmt.Errorf("user already exists")
			}
			return err
		}

		profile := Profile{
			UserID: user.ID,
		}
		if err := tx.Create(&profile).Error; err != nil {
			s.StatusCode = http.StatusInternalServerError
			logger.Error(err.Error())
			return err
		}

		user.ProfileID = profile.ID
		if err := tx.Save(&user).Error; err != nil {
			logger.Error(err.Error())
			return err
		}

		s.UserID = user.ID
		s.Email = user.Email
		s.UserName = user.UserName
		s.User = &user

		return arg.AfterCreate()
	})
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("transaction error: %v", err)
	}

	return nil
}

func (s *UserSerializer) FindCurrentUser(ctx *gin.Context, db *gorm.DB, UserID uint) error {

	logger := common.NewLogger()

	payload, err := GetAuthPayloadFromCtx(ctx)
	if err != nil {
		s.StatusCode = http.StatusUnauthorized
		logger.Error(err.Error())
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}

	var user User

	err = db.Preload("Profile").Where("id = ?", payload.UserID).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			s.StatusCode = http.StatusNotFound
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}

	s.User = &user
	s.UserName = user.UserName
	s.UserID = user.ID
	s.Email = user.Email

	return nil
}

func (s *UserSerializer) FindByEmail(ctx *gin.Context, db *gorm.DB, email string) error {
	var user User

	err := db.Model(&User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		s.StatusCode = http.StatusNotFound
		return fmt.Errorf("unable to retrieve user with email %v %w", s.Email, err)
	}

	s.User = &user
	s.UserName = user.UserName
	s.UserID = user.ID

	return nil
}

func (s *UserSerializer) FindByUID(ctx *gin.Context, db *gorm.DB, uid string) (User, error) {
	logger := common.NewLogger()

	var user User

	err := db.Model(&User{}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		logger.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.StatusCode = http.StatusNotFound
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("unable to retrieve user with uid %v %w", uid, err)
	}

	s.User = &user

	return user, nil
}

func (s *UserSerializer) Search(ctx *gin.Context, db *gorm.DB) ([]User, error) {

	db = buildSearchQuery(db, User{
		UserName: s.UserName,
		Email:    s.Email,
	})
	var users []User

	if err := db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("unable to retrieve users %w", err)
	}

	return users, nil

}

func (s *UserSerializer) Update(ctx *gin.Context, db *gorm.DB) error {
	logger := common.NewLogger()

	payload, err := GetAuthPayloadFromCtx(ctx)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}

	var user User

	err = db.Model(&User{}).Where("id = ?", payload.UserID).First(&user).Error
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}

	user.UserName = s.UserName
	user.Email = s.Email

	if err := db.Model(&User{}).Save(&user).Error; err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to update user %w", err)
	}

	s.User = &user

	return nil
}

func (s *UserSerializer) Delete(ctx *gin.Context, db *gorm.DB) error {
	logger := common.NewLogger()

	payload, err := GetAuthPayloadFromCtx(ctx)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}

	var user User

	err = db.Model(&User{}).Where("id = ?", payload.UserID).First(&user).Error
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}

	if err := db.Model(&User{}).Delete(&user).Error; err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("unable to delete user %w", err)
	}

	return nil
}

func (s *UserSerializer) List(ctx *gin.Context, db *gorm.DB) ([]User, error) {
	logger := common.NewLogger()
	logger.Info(s.Page, s.Limit)
	offset := (s.Page - 1) * s.Limit

	var users []User
	err := db.Limit(s.Limit).Offset(offset).Find(&users).Error
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("unable to retrieve users %w", err)
	}

	if len(users) == 0 {
		s.StatusCode = http.StatusNotFound
		return nil, fmt.Errorf("users not found")
	}

	return users, nil
}

type Preference struct {
	LocationProximity   *string   `json:"locationProximity"`
	Amenities           *[]string `json:"amenities"`
	PriceRange          *[]int    `json:"priceRange"`
	TypeOfAccommodation *string   `json:"typeOfAccommodation"`
}

func (s *UserSerializer) GetPreferences(ctx *gin.Context, db *gorm.DB, uid string) (Preference, error) {
	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return Preference{}, err
	}

	var prefs Preference

	if strings.TrimSpace(user.Preferences) == "{}" {
		// if it's an empty JSON object, just initialize the slice and skip unmarshaling
		prefs = Preference{}
	} else {
		err = json.Unmarshal([]byte(user.Preferences), &prefs)
		if err != nil {
			return prefs, err
		}
	}
	return prefs, nil

}

func (s *UserSerializer) UpdatePreferences(ctx *gin.Context, db *gorm.DB, uid string) error {
	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s.Preferences)
	if err != nil {
		return err
	}

	user.Preferences = string(data)

	if err := db.Model(&User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}

func (s *UserSerializer) DeletePreferences(ctx *gin.Context, db *gorm.DB, uid string) error {
	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return err
	}

	// Clear the Preferences field
	user.Preferences = "{}" // Set it to an empty JSON object to stay consistent

	// Update only the Preferences field in the database
	if err := db.Model(&User{}).Where("id = ?", user.ID).Update("preferences", user.Preferences).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}

type PastSearch map[string]interface{}

func (s *UserSerializer) GetPastSearches(ctx *gin.Context, db *gorm.DB, uid string) ([]PastSearch, error) {

	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return nil, err
	}

	var pastSearches []PastSearch

	if strings.TrimSpace(user.PastSearches) == "{}" {
		// if it's an empty JSON object, just initialize the slice and skip unmarshaling
		pastSearches = []PastSearch{}
	} else {
		err = json.Unmarshal([]byte(user.PastSearches), &pastSearches)
		if err != nil {
			return nil, err
		}
	}

	return pastSearches, nil
}
func (s *UserSerializer) AddPastSearch(ctx *gin.Context, db *gorm.DB, uid string) error {
	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return err
	}

	// Here, instead of overwriting, you'd typically want to append to the existing list.
	var existingSearches []PastSearch

	if strings.TrimSpace(user.PastSearches) == "{}" {
		// if it's an empty JSON object, just initialize the slice and skip unmarshaling
		existingSearches = []PastSearch{}
	} else {
		err = json.Unmarshal([]byte(user.PastSearches), &existingSearches)
		if err != nil {
			return err
		}
	}

	existingSearches = append(existingSearches, s.PastSearches...)

	// Now we marshal the appended searches
	data, err := json.Marshal(existingSearches)
	if err != nil {
		return err
	}

	user.PastSearches = string(data)

	if err := db.Model(&User{}).Where("id = ?", user.ID).Update("past_searches", user.PastSearches).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}

func (s *UserSerializer) ClearPastSearches(ctx *gin.Context, db *gorm.DB, uid string) error {
	user, err := s.FindByUID(ctx, db, uid)
	if err != nil {
		return err
	}

	// Clear the PastSearches field
	user.PastSearches = "{}" // Set it to an empty JSON object to stay consistent

	// Update only the PastSearches field in the database
	if err := db.Model(&User{}).Where("id = ?", user.ID).Update("past_searches", user.PastSearches).Error; err != nil {
		s.StatusCode = http.StatusInternalServerError
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}
