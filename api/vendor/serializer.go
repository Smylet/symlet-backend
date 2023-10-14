package vendor

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/sms"
)

type VendorSerializer struct {
	CompanyName string `json:"company_name"  custom_binding:requiredForCreate`
	Address     string `json:"address" custom_binding:requiredForCreate`

	Email       string `json:"email" custom_binding:requiredForCreate`
	Phone       string `json:"phone" custom_binding:requiredForCreate`
	Website     string `json:"website" custom_binding:requiredForCreate`
	Description string `json:"description" custom_binding:requiredForCreate`
	Service     string `json:"service" custom_binding:requiredForCreate`
	Vendor      *Vendor `json:"-`
}

func (s *VendorSerializer) Create(ctx *gin.Context, db *gorm.DB, sms *sms.SMSSender) error {
	s.Vendor = &Vendor{
		CompanyName: s.CompanyName,
		Address:     s.Address,
		Email:       s.Email,
		Phone:       s.Phone,
		Website:     s.Website,
		Description: s.Description,
		Service:     s.Service,
	}

	payload, err := users.GetAuthPayloadFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}
	// Does this User already have a Profile?
	fmt.Printf("payload %v", payload)
	err = db.Model(&users.User{}).Preload(clause.Associations).Where("id = ?", payload.UserID).First(&s.Vendor.User).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}
	if s.Vendor.User.RoleType != ""{
		return fmt.Errorf(
			"user is as already associated with a %v role", s.Vendor.User.RoleType)
	}
	s.Vendor = &Vendor{
		User: s.Vendor.User,
		CompanyName: s.CompanyName,
		Address:     s.Address,
		Email:       s.Email,
		Phone:       s.Phone,
		Website:     s.Website,
		Description: s.Description,
		Service:     s.Service,
	}
 	err = common.ExecTx(ctx, db, func(tx *gorm.DB) error {
	if err = db.Create(&s.Vendor).Error; err != nil {
			return fmt.Errorf("unable to create vendor %w", err)
		}
		if err = sms.SendSMS(s.Vendor.Phone, "Welcome to Symlet, your account has been created"); err!= nil {
				return fmt.Errorf("unable to send sms %w", err)
			
		}
		return nil
	},
	)
	if err != nil {
		return fmt.Errorf("unable to create vendor %w", err)
	}

	return nil
}

func (s *VendorSerializer)Get(c *gin.Context, db *gorm.DB, uid string) error {
	err := db.Model(&Vendor{}).Preload(clause.Associations).Where("uid = ?", uid).First(&s.Vendor).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve vendor with uid of %v, %w", uid, err)
	}

	return nil
}

func (s *VendorSerializer) Response() map[string]interface{} {
	return map[string]interface{}{
		"uid":          s.Vendor.UID,
		"id":           s.Vendor.ID,
		"company_name": s.Vendor.CompanyName,
		"address":      s.Vendor.Address,
		"email":        s.Vendor.Email,
		"phone":        s.Vendor.Phone,
		"website":      s.Vendor.Website,
		"description":  s.Vendor.Description,
		"service":      s.Vendor.Service,
		"created_at":   s.Vendor.CreatedAt,
		"updated_at":   s.Vendor.UpdatedAt,
		"is_verified":  s.Vendor.IsVerified,
		"user": map[string]interface{}{
			"uid":        s.Vendor.User.UID,
			"username":   s.Vendor.User.Username,
			"email":      s.Vendor.User.Email,
			"created_at": s.Vendor.User.CreatedAt,
			"updated_at": s.Vendor.User.UpdatedAt,
		},
		"profile": map[string]interface{}{
			"uid":        s.Vendor.User.Profile.UID,
			"first_name": s.Vendor.User.Profile.FirstName,
			"last_name":  s.Vendor.User.Profile.LastName,
			"bio":        s.Vendor.User.Profile.Bio,
			"image":      s.Vendor.User.Profile.Image,
			"created_at": s.Vendor.User.Profile.CreatedAt,
			"updated_at": s.Vendor.User.Profile.UpdatedAt,
		},
	}

}
