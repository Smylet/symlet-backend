package hostel

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/utilities/common"
	//"github.com/Smylet/symlet-backend/utilities/db"
)



type Map map[string]float64

type Hostel struct {
	common.AbstractBaseModel
	Name         string                        `gorm:"not null"`
	UniversityID uint                          `gorm:"not null"`
	University   reference.ReferenceUniversity `gorm:"foreignKey:UniversityID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Address      string                        `gorm:"not null"`
	City         string                        `gorm:"not null"`
	State        string                        `gorm:"not null"`
	Country      string                        `gorm:"not null"`
	Description  string                        `gorm:"not null"`

	ManagerID uint                  `gorm:"not null"`
	Manager   manager.HostelManager `gorm:"foreignKey:ManagerID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	Amenities []*reference.ReferenceHostelAmenities `gorm:"many2many:hostel_ammenities;constraint:OnDelete:CASCADE"`
	Students  []*student.Student                    `gorm:"many2many:hostel_students;"`

	// Other features
	NumberOfUnits         uint   `gorm:"not null"`
	NumberOfOccupiedUnits uint   `gorm:"not null"`
	NumberOfBedrooms      uint   `gorm:"not null"`
	NumberOfBathrooms     uint   `gorm:"not null"`
	Kitchen               string `gorm:"not null; oneof=shared none private"`
	FloorSpace            uint   `gorm:"not null"`
	HostelFee             HostelFee `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	HostelImages          []HostelImage `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type HostelImage struct {
	common.AbstractBaseImageModel
	HostelID uint
	Hostel   Hostel `gorm:"foreignKey:HostelID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type HostelFee struct {
	common.AbstractBaseModel
	HostelID    uint 
	TotalAmount float64
	PaymentPlan string `gorm:"oneof: 'monthly' 'by_school_session' 'annually'"`
	Breakdown   Map    `gorm:"type:jsonb"`
}

type HostelAgreementTemplate struct {
	common.AbstractBaseModel
	HostelID    uint   `gorm:"not null"`
	DocumentURL string `gorm:"not null"`
}

func (m Map) Value() (driver.Value, error) {
	byte, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(byte), nil
}

func (m *Map) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	switch data := v.(type) {
	case string:
		return json.Unmarshal([]byte(data), &m)
	case []byte:
		return json.Unmarshal(data, &m)
	default:
		return fmt.Errorf("cannot scan type %t into Map", v)
	}
}
