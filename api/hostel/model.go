package hostel

import (
	"time"

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
	Name       string             `gorm:"not null"`
	UniversityID uint            `gorm:"not null"`
	University reference.ReferenceUniversity   `gorm:"foreignKey:UniversityID"`
	Address    string             `gorm:"not null"`
	City       string             `gorm:"not null"`
	State      string             `gorm:"not null"`
	Country    string             `gorm:"not null"`
	Description string            `gorm:"not null"`


	ManagerID  uint               `gorm:"not null"`
	Manager    users.User         `gorm:"foreignKey:ManagerID"`

	Ammenities []*reference.ReferenceHostelAmmenities `gorm:"many2many:hostel_ammenities;"`
	//Students   []*student.Student `gorm:"many2many:hostel_students;"`
	
	// Other features
	NumberOfUnits uint            `gorm:"not null"`
	NumberOfOccupiedUnits uint    `gorm:"not null"`
	NumberOfBedrooms uint         `gorm:"not null"`
	NumberOfBathrooms uint        `gorm:"not null"`
	Kitchen    bool               `gorm:"not null"`
	FloorSpace uint              `gorm:"not null"`
}

type HostelImage struct{
	core.AbstractBaseModel
	ImageURLs []string
	HostelID uint
	Hostel Hostel `gorm:"foreignKey:HostelID"`
}

// HostelStudent is the join table between Hostel and Student
// It holds the relationship between the two entities and other metadata
type HostelStudent struct {
	core.AbstractBaseModel
	StudentID     uint
	HostelID      uint
	CheckInDate   time.Time
	CheckOutDate  time.Time
	RoomNumber    string
	CurrentHostel bool
	Student       student.Student `gorm:"foreignKey:StudentID"`
	Hostel        Hostel     `gorm:"foreignKey:HostelID"`
	
	SignedAgreement bool
	HostelAgreementTemplateID uint 
	HostelAgreementTemplate HostelAgreementTemplate `gorm:"foreignKey:HostelAgreementTemplateID"`
	SubmittedSignedAgreementUrl string

	// Other metadata fields specific to the student-hostel relationship
}


type HostelFee struct {
	core.AbstractBaseModel
	HostelID uint 
	Hostel Hostel
	TotalAmount   float64
	Breakdown     []HostelFeeBreakdown `gorm:"embedded"`
}

type HostelFeeBreakdown struct {
	Description string  `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
}

type HostelAgreementTemplate struct {
	core.AbstractBaseModel 
	HostelID uint `gorm:"not null"`
	Hostel Hostel 
	DocumentURL string `gorm:"not null"`
}
