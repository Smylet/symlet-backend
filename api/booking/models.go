package booking

import (
	"database/sql"
	"time"

	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/utilities/common"
)

// HostelStudent is the join table between Hostel and Student
// It holds the relationship between the two entities and other metadata
type HostelStudent struct {
	common.AbstractBaseModel
	StudentID     uint `gorm:"primaryKey"`
	HostelID      uint `gorm:"primaryKey"`
	CheckInDate   time.Time
	CheckOutDate  sql.NullTime
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
	LastPaymentDate sql.NullTime
	NextPaymentDate sql.NullTime
	HostelBookingID uint
	HostelBooking   HostelBooking `gorm:"foreignKey:HostelBookingID;constraint:OnDelete:CASCADE"`
}

// Booking model
type HostelBooking struct {
	common.AbstractBaseModel
	StudentID    uint `gorm:"not null"`
	Student      student.Student
	HostelID     uint `gorm:"not null"`
	Hostel       hostel.Hostel
	PartStay     bool `gorm:"default:false"`
	NumberOfMonths int `gorm:"default:12"`
	PaymentPlans []PaymentPlan `gorm:"foreignKey:HostelBookingID;constraint:OnDelete:CASCADE"`
}

type PaymentPlan struct {
	common.AbstractBaseModel
	Amount               float64 `gorm:"not null"`
	HostelBookingID      uint    `gorm:"not null"`
	HostelBooking        HostelBooking
	// allowed values: all, spread, stay, deferred
	// all: pay all at once
	// spread: spread the payment across multiple months
	// deferred: pay at a later date
	PaymentType          string         `gorm:"not null;default:'all';check:payment_type IN ('all', 'spread', 'deferred')"`
	PaymentInterval      sql.NullString `gorm:"check:payment_interval IN ('equal', 'unequal')"` // not needed if ALL
	IntervalDuration     sql.NullInt32  // `gorm:"check:interval_duration LESS"` // Only for 'equal' distribution
	DeferredDate         sql.NullTime   // Only for 'deferred' payment
	PaymentDistributions []PaymentDistribution
}

// PaymentDistribution model
type PaymentDistribution struct {
	common.AbstractBaseModel
	PaymentPlanID uint      `gorm:"not null"`
	Date          time.Time `gorm:"not null"`
	Amount        float64   `gorm:"not null"`
}
