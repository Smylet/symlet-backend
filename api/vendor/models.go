package vendor

import (
	"github.com/Smylet/symlet-backend/api/core"
)

type Vendor struct {
	core.AbstractBaseModel
	UserID uint `gorm:"not null"`

	CompanyName string `gorm:"not null"`
	Address string `gorm:"not null"`
	Email string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Website string `gorm:"not null"`
	Logo string `gorm:"not null"`
	Description string `gorm:"not null"`
	Service string `gorm:"not null"`
	Rating float64 `gorm:"not null"`
	IsVerified bool `gorm:"default:false"`
}