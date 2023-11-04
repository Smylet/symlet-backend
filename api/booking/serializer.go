package booking

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
)

type PaymentDistributionSerializer struct {
	PaymentPlanID uint      `json:"payment_plan_id"`
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
}

type PaymentPlanSerializer struct {
	Amount               float64 `json:"amount" binding:"required"`
	HostelBookingID      uint    `json:"-"`
	PaymentType          string  `json:"payment_type" binding:"required,oneof=deferred all spread"`
	PaymentInterval      string  `json:"payment_interval" binding:"required,oneof=equal unequal"`
	IntervalDuration     int32   `json:"interval_duration"`
	DeferredDate         *time.Time
	PaymentDistributions []PaymentDistributionSerializer `json:"payment_distributions"`
}

type BookingSerializer struct {
	HostelUID      string                `json:"hostel_uid"`
	StudentUID     string                `json:"student_uid"`
	PartStay       bool                  `json:"part_stay"`
	NumberOfMonths uint32                `json:"number_of_months"` // Only for 'stay' payment
	PaymentPlan    PaymentPlanSerializer `json:"payment_plan"`
	HostelBooking 	  *HostelBooking	`json:"-"`
}

func (b *BookingSerializer) Validate() error {
	payment_plan := b.PaymentPlan
	if err := VerifyPaymentPlanSerializer(payment_plan); err != nil {
		return err
	}

	return nil
}

func (b *BookingSerializer) Create(ctx *gin.Context, db *gorm.DB) error {
	hostelUID := ctx.Param("uid")
	hst := &hostel.Hostel{}
	err := db.Model(&hostel.Hostel{}).Where("uid = ?", hostelUID).First(hst).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve hostel with uid %v %w", hostelUID, err)
	}
	b.HostelUID = hostelUID 
	user, err := users.GetAuthPayloadFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}
	std := &student.Student{}
	err = db.Model(&student.Student{}).Where("user_id = ?", user.UserID).First(std).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve student with user id %v %w", user.UserID, err)
	}
	err = b.Validate()
	if err != nil {
		return err
	}
	b.HostelBooking = &HostelBooking{
		StudentID:    std.ID,
		HostelID:     hst.ID,
		PartStay:     b.PartStay,
		NumberOfMonths:   int(b.NumberOfMonths),
	}
	err = common.ExecTx(ctx, db, func(tx *gorm.DB)error{
		err := tx.Create(&b.HostelBooking).Error
		if err != nil {
			return fmt.Errorf("unable to create hostel booking %w", err)
		}
		paymentPlan := &PaymentPlan{
			Amount:               b.PaymentPlan.Amount,
			HostelBookingID:      b.HostelBooking.ID,
			PaymentType:          b.PaymentPlan.PaymentType,
			PaymentInterval:      sql.NullString{
				String: b.PaymentPlan.PaymentInterval,
				Valid:  true,
			},//b.PaymentPlan.PaymentInterval,
			IntervalDuration:     sql.NullInt32{
				Int32: b.PaymentPlan.IntervalDuration,
				Valid: true,
			},//b.PaymentPlan.IntervalDuration,
			DeferredDate:         sql.NullTime{
				Time:  *b.PaymentPlan.DeferredDate,
				Valid: b.PaymentPlan.DeferredDate != nil,
			},//b.PaymentPlan.DeferredDate,
		}
		err = tx.Create(&paymentPlan).Error
		if err != nil {
			return fmt.Errorf("unable to create payment plan %w", err)
		}
		for _, paymentDistribution := range b.PaymentPlan.PaymentDistributions {
			paymentDistribution := &PaymentDistribution{
				PaymentPlanID: paymentPlan.ID,
				Date:          paymentDistribution.Date,
				Amount:        paymentDistribution.Amount,
			}
			err = tx.Create(&paymentDistribution).Error
			if err != nil {
				return fmt.Errorf("unable to create payment distribution %w", err)
			}
		}
		return nil

	})
	if err != nil {
		return fmt.Errorf("unable to create hostel booking %w", err)
	}
	return nil
}


func VerifyPaymentPlanSerializer(ser PaymentPlanSerializer) error {
	if ser.PaymentType == "deferred" {
		return verifyDeferredPayment(ser)
	} else if ser.PaymentType == "all" {
		return verifyAllPayment(ser)
	} else if ser.PaymentType == "spread" {
		return verifySpreadPayment(ser)
	} else {
		return errors.New("payment_type is invalid")
	}
}

func verifyDeferredPayment(ser PaymentPlanSerializer) error {
	if ser.DeferredDate.IsZero() {
		return errors.New("deferred_date is required for deferred payment")
	}
	ser.PaymentInterval = "equal"
	ser.IntervalDuration = 1
	ser.PaymentDistributions = []PaymentDistributionSerializer{
		{
			//PaymentPlanID: ser.HostelBookingID,
			Date:   *ser.DeferredDate,
			Amount: ser.Amount,
		},
	}
	return nil
}

func verifyAllPayment(ser PaymentPlanSerializer) error {
	ser.DeferredDate = nil
	ser.PaymentDistributions = []PaymentDistributionSerializer{
		{
			//PaymentPlanID: ser.HostelBookingID,
			Date:   time.Now(),
			Amount: ser.Amount,
		},
	}
	ser.IntervalDuration = 1
	ser.PaymentInterval = "equal"
	return nil
}

func verifySpreadPayment(ser PaymentPlanSerializer) error {
	if ser.PaymentInterval == "" {
		return errors.New("payment_interval is required for spread payment")
	}
	if ser.PaymentInterval == "equal" {
		if ser.IntervalDuration == 0 {
			return errors.New("payment_interval is required for equal payment interval")
		} else {
			return verifyEqualSpreadPayment(ser)
		}
	} else if ser.PaymentInterval == "unequal" {
		return verifyUnequalSpreadPayment(ser)
	} else {
		return errors.New("payment_type is invalid")
	}
}

func verifyEqualSpreadPayment(ser PaymentPlanSerializer) error {
	if ser.IntervalDuration <= 0 {
		return errors.New("interval duration must be greater than zero for equal payment interval")
	}

	amount := 1000.0 // Replace with the actual amount calculation logic

	amountPerStep := amount / float64(ser.IntervalDuration)
	nextDate := time.Now()

	for i := 0; i < int(ser.IntervalDuration); i++ {
		ser.PaymentDistributions = append(ser.PaymentDistributions, PaymentDistributionSerializer{
			//PaymentPlanID: ser.HostelBookingID,
			Date:   nextDate,
			Amount: amountPerStep,
		})
		nextDate = nextDate.AddDate(0, 0, 1) // Add one day for equal intervals
	}

	ser.DeferredDate = nil
	return nil
}

func verifyUnequalSpreadPayment(ser PaymentPlanSerializer) error {
	if len(ser.PaymentDistributions) == 0 {
		return errors.New("payment_distributions is required for unequal spread payment")
	}

	totalAmount := 0.0
	today := time.Now()

	for _, dist := range ser.PaymentDistributions {
		amount := dist.Amount
		if amount == 0.0 {
			return errors.New("amount is required for unequal spread payment")
		}
		date := dist.Date
		if date.IsZero() {
			return errors.New("date is required for unequal spread payment")
		}

		if date.Before(today) {
			return errors.New("date cannot be in the past")
		}
		// TODO: Add validation for date range when it is a part stay
		if date.After(today.AddDate(0, 12, 0)) {
			return errors.New("date cannot be more than 12 months in the future")
		}

		totalAmount += amount
	}

	// Replace this with your actual amount calculation logic
	actualAmount := 1000.0

	if totalAmount != actualAmount {
		return errors.New("total amount entered does not match the actual total amount")
	}

	return nil
}
