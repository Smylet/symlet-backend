package booking

import (
	"database/sql"
	"time"

	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
)



// HostelStudent is the join table between Hostel and Student
// It holds the relationship between the two entities and other metadata
type HostelStudent struct {
	core.AbstractBaseModel
	StudentID     uint `gorm:"primaryKey"`
	HostelID      uint	`gorm:"primaryKey"`
	CheckInDate   time.Time
	CheckOutDate  time.Time
	RoomNumber    *string
	CurrentHostel bool
	
	// Agreement Related
	SignedAgreement             bool `gorm:"default:false"`
	HostelAgreementTemplateID   sql.NullInt64
	HostelAgreementTemplate     *hostel.HostelAgreementTemplate `gorm:"foreignKey:HostelAgreementTemplateID;constraint:OnDelete:SET NULL"`
	SubmittedSignedAgreementUrl string
	
	// Payment Related
	TotalAmountPaid float64 `gorm:"default:0"`
	TotalAmountDue  float64 `gorm:"default:0"`
	LastPaymentDate time.Time
	NextPaymentDate time.Time
	HostelBookingID uint
	HostelBooking HostelBooking `gorm:"foreignKey:HostelBookingID;constraint:OnDelete:CASCADE"`
}

// Booking model
type HostelBooking struct {
	core.AbstractBaseModel
	StudentID uint `gorm:"not null"`
	Student student.Student 
	HostelID uint `gorm:"not null"`
	Hostel hostel.Hostel
	CustomPaymentPlans []CustomPaymentPlan `gorm:"foreignKey:HostelBookingID;constraint:OnDelete:CASCADE"`

}
	

type CustomPaymentPlan struct {
	core.AbstractBaseModel
	HostelBookingID uint `gorm:"not null"`
	HostelBooking HostelBooking
	PaymentType      string           `gorm:"not null;check:payment_type IN ('spread', 'stay', 'deferred')"`
	PaymentInterval  string           `gorm:"not null;check:payment_interval IN ('equal', 'unequal')"`
	IntervalDuration string           `gorm:"check:interval_duration IN ('3', '6')"` // Only for 'equal' distribution
	DeferredDate     *time.Time       // Only for 'deferred' payment
	NumberOfMonths   sql.NullInt32             // Only for 'stay' payment
	PaymentDistributions []PaymentDistribution
}

// PaymentDistribution model
type PaymentDistribution struct {
	core.AbstractBaseModel
	CustomPaymentPlanID uint             `gorm:"not null"`
	Date 			  time.Time        `gorm:"not null"`
	Amount              float64          `gorm:"not null"`
}
