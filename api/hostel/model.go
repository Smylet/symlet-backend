package hostel

import (

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"

	"github.com/Smylet/symlet-backend/api/core"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Hostel{},
// 		&HostelStudent{},
// 	)
// }

type Hostel struct {
	core.AbstractBaseModel
	Name         string                        `gorm:"not null"`
	UniversityID uint                          `gorm:"not null"`
	University   reference.ReferenceUniversity `gorm:"foreignKey:UniversityID;constraint:OnDelete:SET NULL"`
	Address      string                        `gorm:"not null"`
	City         string                        `gorm:"not null"`
	State        string                        `gorm:"not null"`
	Country      string                        `gorm:"not null"`
	Description  string                        `gorm:"not null"`

	ManagerID uint       `gorm:"not null"`
	Manager   users.User `gorm:"foreignKey:ManagerID"`

	Ammenities []*reference.ReferenceHostelAmmenities `gorm:"many2many:hostel_ammenities;"`
	Students   []*student.Student `gorm:"many2many:hostel_students;"`

	// Other features
	NumberOfUnits         uint `gorm:"not null"`
	NumberOfOccupiedUnits uint `gorm:"not null"`
	NumberOfBedrooms      uint `gorm:"not null"`
	NumberOfBathrooms     uint `gorm:"not null"`
	Kitchen               bool `gorm:"not null"`
	FloorSpace            uint `gorm:"not null"`
	HostelFee             HostelFee
	HostelImages          []HostelImage `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE"`
}

type HostelImage struct {
	core.AbstractBaseImageModel
	HostelID  uint
	Hostel    Hostel `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE"`
}


type HostelFee struct {
	core.AbstractBaseModel
	HostelID    uint
	TotalAmount float64
	PaymentPlan string `gorm:"oneof: 'monthly' 'by_school_session' 'annually'"`
	Breakdown   map[string]interface{} `gorm:"type:json"`

}


type HostelAgreementTemplate struct {
	core.AbstractBaseModel
	HostelID    uint   `gorm:"not null"`
	DocumentURL string `gorm:"not null"`
}
