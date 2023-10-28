package users

import (
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileSerializer struct {
	FirstName string                `json:"first_name" form:"first_name" custom_binding:"requiredForCreate"`
	LastName  string                `json:"last_name" form:"last_name" custom_binding:"requiredForCreate"`
	Bio       string                `json:"bio" form:"bio"`
	Image     *multipart.FileHeader `form:"image" swaggerignore:"true"`
	Profile   *Profile              `json:"-"`
}

// func (s *ProfileSerializer) Update(ctx *gin.Context, db *gorm.DB, AWSsession *session.Session) error {
// 	payload, err := GetAuthPayloadFromCtx(ctx)
// 	if err != nil {
// 		return fmt.Errorf("unable to retrieve user payload from context %w", err)
// 	}
// 	// Does this User already have a Profile?
// 	err = db.Model(&Profile{}).Where("user_id = ?", payload.UserID).First(&s.Profile).Error

// 	if err == nil {
// 		return fmt.Errorf("User does not have a profile, %w", err)
// 	}

// 	s.Profile = &Profile{
// 		UserID:    payload.UserID,
// 		FirstName: s.FirstName,
// 		LastName:  s.LastName,
// 		Bio:       s.Bio,
// 	}
// 	filePath, err := utils.ProcessUploadedImage(s.Image, AWSsession)
// 	if err != nil{
// 		return fmt.Errorf("unable to upload images %w", err)
// 	}
// 	s.Profile.Image = filePath
// 	err = db.Model(&Profile{}).Create(s.Profile).Error
// 	if err != nil{
// 		return fmt.Errorf("unable to create Profile, %w", err)
// 	}
// 	return nil

// }
func (s *ProfileSerializer) Get(ctx *gin.Context, db *gorm.DB, uid string) error {
	// Does this User already have a Profile?
	err := db.Model(&Profile{}).Where("uid = ?", uid).First(&s.Profile).Error
	if err != nil {
		return fmt.Errorf("User does not have a profile")
	}

	return nil
}

func (s *ProfileSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":         s.Profile.ID,
		"uid":        s.Profile.UID,
		"created_at": s.Profile.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"updated_at": s.Profile.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		"bio":        s.Profile.Bio,
		"image":      s.Profile.Image,
		"first_name": s.Profile.FirstName,
		"last_name":  s.Profile.LastName,
		"user_id":    s.Profile.UserID,
	}
	return response
}

type VerifyEmailRequest struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // This can be used to send the expiration time of the token, if needed.
}
