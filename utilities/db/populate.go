package db

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/booking"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/api/vendor"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/go-faker/faker/v4"
)

func createVendor(ctx context.Context, db *gorm.DB) (vendor.Vendor, error) {
	vendorUser := users.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	vendor := vendor.Vendor{
		CompanyName: faker.Word(),
		Address:     faker.Sentence(),
		Email:       faker.Email(),
		Phone:       faker.Phonenumber(),
		Website:     faker.URL(),
		Logo:        faker.URL(),
		Description: faker.Sentence(),
		Service:     faker.Sentence(),
		IsVerified:  true,
	}

	err := common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		if err := tx.Create(&vendorUser).Error; err != nil {
			return err
		}
		vendor.User = vendorUser
		if err := tx.Create(&vendor).Error; err != nil {
			return err
		}
		return nil
	})

	return vendor, err
}

func createHostelManager(ctx context.Context, db *gorm.DB) (manager.HostelManager, error) {
	managerUser := users.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	hostelManager := manager.HostelManager{}

	err := common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		if err := tx.Create(&managerUser).Error; err != nil {
			return err
		}
		hostelManager.User = managerUser
		if err := tx.Create(&hostelManager).Error; err != nil {
			return err
		}
		return nil
	})

	return hostelManager, err
}

func createHostel(ctx context.Context, db *gorm.DB, university *reference.ReferenceUniversity, ammenities []*reference.ReferenceHostelAmenities) (hostel.Hostel, error) {
	hostelManager, err := createHostelManager(ctx, db)
	if err != nil {
		return hostel.Hostel{}, err
	}

	hostelObj := hostel.Hostel{
		Manager:               hostelManager,
		University:            *university,
		Name:                  faker.Word(),
		Address:               faker.Sentence(),
		City:                  "Ilorin",
		State:                 "Kwara",
		Country:               "Nigeria",
		Description:           faker.Paragraph(),
		NumberOfUnits:         10,
		NumberOfOccupiedUnits: 0,
		NumberOfBedrooms:      1,
		NumberOfBathrooms:     1,
		Kitchen:               "detached",
		FloorSpace:            100,
		HostelFee: hostel.HostelFee{
			TotalAmount: 100000,
			PaymentPlan: "monthly",
			Breakdown: map[string]float64{
				"rent":           100000,
				"service_charge": 0,
				"caution_fee":    0,
				"agency_fee":     0,
			},
		},
		Amenities: ammenities,
	}

	err = db.Create(&hostelObj).Error
	if err != nil {
		return hostel.Hostel{}, err
	}

	return hostelObj, nil
}

func createStudent(ctx context.Context, db *gorm.DB, university reference.ReferenceUniversity) (student.Student, error) {
	studentUser := users.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	studentObj := student.Student{
		University: university,
	}

	err := common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		if err := tx.Create(&studentUser).Error; err != nil {
			return err
		}
		studentObj.User = studentUser
		if err := tx.Create(&studentObj).Error; err != nil {

			return err
		}
		return nil
	})
	if err != nil {
		return student.Student{}, err
	}

	// Populate student data and create records

	return studentObj, nil
}

func createHostelBooking(ctx context.Context, db *gorm.DB, hostel hostel.Hostel, student student.Student) (booking.HostelBooking, error) {
	// Create and return a hostel booking record

	hostelBooking := booking.HostelBooking{
		Hostel:  hostel,
		Student: student,
		PaymentPlans: []booking.PaymentPlan{
			{
				Amount:      100000,
				PaymentType: "all",
				PaymentInterval: sql.NullString{
					String: "equal",
					Valid:  true,
				},
			},
		},
	}
	err := db.Create(&hostelBooking).Error
	if err != nil {
		return booking.HostelBooking{}, err
	}
	return hostelBooking, nil
}

func createHostelStudent(ctx context.Context, db *gorm.DB, hostel hostel.Hostel, student student.Student, hostelBooking booking.HostelBooking) (booking.HostelStudent, error) {
	// Create and return a hostel student record
	randomDate, err := time.Parse("2006/01/02", faker.Date())
	if err != nil {
		return booking.HostelStudent{}, err
	}
	roomNumber := "45A"

	hostel_student := booking.HostelStudent{
		HostelID:        hostel.ID,
		StudentID:       student.ID,
		HostelBooking:   hostelBooking,
		CheckInDate:     randomDate,
		RoomNumber:      &roomNumber,
		CurrentHostel:   true,
		TotalAmountPaid: 100000,
		TotalAmountDue:  0,
		SignedAgreement: false,
		// HostelAgreementTemplate: hostel.HostelAgreementTemplate{
		// 	DocumentURL: faker.URL(),

		// },
		// SubmittedSignedAgreementUrl: faker.URL(),
	}
	hostel_student.CheckOutDate = sql.NullTime{
		Time:  hostel_student.CheckInDate.AddDate(1, 0, 0),
		Valid: true,
	}

	hostel_student.NextPaymentDate = sql.NullTime{
		Time:  hostel_student.CheckInDate.AddDate(1, 0, 0),
		Valid: true,
	}

	err = db.Create(&hostel_student).Error
	if err != nil {
		return booking.HostelStudent{}, err
	}

	return booking.HostelStudent{}, nil
}

func PopulateDatabase(db *gorm.DB) error {
	// Migrate DB With Atlas

	ctx := context.Background()

	university := reference.ReferenceUniversity{}
	ammenities := []*reference.ReferenceHostelAmenities{}

	err := common.ExecTx(ctx, db, func(tx *gorm.DB) error {
		err := db.Model(&university).Limit(1).First(&university).Error
		if err != nil {
			return err
		}

		err = db.Model(&ammenities).Limit(10).Find(&ammenities).Error
		if err != nil {
			return err
		}

		_, err = createVendor(ctx, db)
		if err != nil {
			return err
		}

		hostel, err := createHostel(ctx, db, &university, ammenities)
		if err != nil {
			return err
		}

		student, err := createStudent(ctx, db, university)
		if err != nil {
			return err
		}

		hostelBooking, err := createHostelBooking(ctx, db, hostel, student)
		if err != nil {
			return err
		}

		_, err = createHostelStudent(ctx, db, hostel, student, hostelBooking)
		if err != nil {
			return err
		}
		return nil
	},
	)
	return err
}
