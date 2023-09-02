package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Smylet/symlet-backend/api/booking"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/api/vendor"
	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func createVendor(db *gorm.DB) (vendor.Vendor, error) {
	vendorUser := users.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		PasswordHash: faker.Password(),
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

	err := db.Transaction(func(tx *gorm.DB) error {
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

func createHostelManager(db *gorm.DB) (manager.HostelManager, error) {
	managerUser := users.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		PasswordHash: faker.Password(),
	}

	hostelManager := manager.HostelManager{
	}

	err := db.Transaction(func(tx *gorm.DB) error {
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

func createHostel(db *gorm.DB, university *reference.ReferenceUniversity, ammenities []*reference.ReferenceHostelAmmenities) (hostel.Hostel, error) {
	hostelManager, err := createHostelManager(db)
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
		Kitchen:               true,
		FloorSpace:            100,
		HostelFee: hostel.HostelFee{
			TotalAmount: 100000,
			PaymentPlan: "monthly",
			Breakdown: map[string]interface{}{
				"rent":           100000,
				"service_charge": 0,
				"caution_fee":    0,
				"agency_fee":     0,
			},
		},
		// HostelImages: []hostel.HostelImage{
		// 	{
		// 		ImageURL: faker.URL(),
		// 	},
		// 	{
		// 		ImageURL: faker.URL(),
		// 	},
		// 	{
		// 		ImageURL: faker.URL(),
		// 	},
		// },
		Ammenities: ammenities,
	}

	err = db.Create(&hostelObj).Error
	if err != nil {
		return hostel.Hostel{}, err
	}

	return hostelObj, nil
}

func createStudent(db *gorm.DB, university reference.ReferenceUniversity) (student.Student, error) {
	studentUser := users.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		PasswordHash: faker.Password(),
	}

	studentObj := student.Student{
		University: university,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
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

func createHostelBooking(db *gorm.DB, hostel hostel.Hostel, student student.Student) (booking.HostelBooking, error) {
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

func createHostelStudent(db *gorm.DB, hostel hostel.Hostel, student student.Student, hostelBooking booking.HostelBooking) (booking.HostelStudent, error) {
	//Create and return a hostel student record
	randomDate, err := time.Parse(time.DateOnly, faker.Date())
	if err != nil {
		return booking.HostelStudent{}, err
	}
	roomNumber := "45A"


	

	hostel_student := booking.HostelStudent{
		HostelID:        hostel.ID,
		StudentID:       student.ID,
		HostelBooking:   hostelBooking,
		CheckInDate:     randomDate,
		RoomNumber:     &roomNumber,
		CurrentHostel:   true,
		TotalAmountPaid: 100000,
		TotalAmountDue:  0,
		SignedAgreement: false,
		// HostelAgreementTemplate: hostel.HostelAgreementTemplate{
		// 	DocumentURL: faker.URL(),

		// },
		//SubmittedSignedAgreementUrl: faker.URL(),
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

func PopulateDatabase(db *gorm.DB) error{
	Migrate(db)

	university := reference.ReferenceUniversity{}
	ammenities := []*reference.ReferenceHostelAmmenities{}

	err := db.Transaction(func(tx *gorm.DB) error {

		err := db.Model(&university).Limit(1).First(&university).Error
		if err != nil {
			return err
		}
		fmt.Print(university)

		err = db.Model(&ammenities).Limit(10).Find(&ammenities).Error
		if err != nil {
			return err
		}


		_, err = createVendor(db)
		if err != nil {
			return err
		}

		hostel, err := createHostel(db, &university, ammenities)
		if err != nil {
			return err
		}

		student, err := createStudent(db, university)
		if err != nil {
			return err

		}

		hostelBooking, err := createHostelBooking(db, hostel, student)
		if err != nil {
			return err
		}

		_, err = createHostelStudent(db, hostel, student, hostelBooking)
		if err != nil {
			return err
		}
		return nil

	},
	)
	return err
}
