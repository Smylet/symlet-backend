package hostel

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/utilities/common"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Hostel{},
// 		&HostelStudent{},
// 	)
// }

type cusjsonb map[string]float64


type Hostel struct {
	common.AbstractBaseModel
	Name         string                        `gorm:"not null"`
	UniversityID uint                          `gorm:"not null"`
	University   reference.ReferenceUniversity `gorm:"foreignKey:UniversityID;constraint:OnDelete:SET NULL"`
	Address      string                        `gorm:"not null"`
	City         string                        `gorm:"not null"`
	State        string                        `gorm:"not null"`
	Country      string                        `gorm:"not null"`
	Description  string                        `gorm:"not null"`

	ManagerID uint                  `gorm:"not null"`
	Manager   manager.HostelManager `gorm:"foreignKey:ManagerID"`

	Amenities []*reference.ReferenceHostelAmenities `gorm:"many2many:hostel_ammenities;"`
	Students  []*student.Student                    `gorm:"many2many:hostel_students;"`

	// Other features
	NumberOfUnits         uint   `gorm:"not null"`
	NumberOfOccupiedUnits uint   `gorm:"not null"`
	NumberOfBedrooms      uint   `gorm:"not null"`
	NumberOfBathrooms     uint   `gorm:"not null"`
	Kitchen               string `gorm:"not null; oneof=shared none private"`
	FloorSpace            uint   `gorm:"not null"`
	HostelFee             HostelFee
	HostelImages          []HostelImage `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE"`
}

type HostelImage struct {
	common.AbstractBaseImageModel
	HostelID uint
	Hostel   Hostel `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE"`
}

type HostelFee struct {
	common.AbstractBaseModel
	HostelID    uint
	TotalAmount float64
	PaymentPlan string                 `gorm:"oneof: 'monthly' 'by_school_session' 'annually'"`
	Breakdown   cusjsonb `gorm:"type:jsonb;not null;default: '{}'::jsonb"`
}

type HostelAgreementTemplate struct {
	common.AbstractBaseModel
	HostelID    uint   `gorm:"not null"`
	DocumentURL string `gorm:"not null"`
}

// Returns the JSON-encoded representation
func (a cusjsonb) Value() (driver.Value, error) {
    // Convert to map[string]float32 from map[int]float32 
    x := make(map[string]float64)

    // Marshal into json 
    return json.Marshal(x)
}

// Decodes a JSON-encoded value
func (a *cusjsonb) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    // Unmarshal from json to map[string]float32
    x := make(map[string]float64)
    if err := json.Unmarshal(b, &x); err != nil {
       return err
    }

    return nil
}