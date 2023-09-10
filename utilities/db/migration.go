package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/api/booking"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/maintenance"
	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/notification"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/review"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/api/vendor"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(

		users.User{},
		users.Profile{},
		users.VerificationEmail{},
		users.Session{},
		student.Student{},

		manager.HostelManager{},
		reference.ReferenceHostelAmenities{},
		reference.ReferenceUniversity{},

		// hostel
		hostel.Hostel{},
		booking.HostelStudent{},
		booking.HostelBooking{},
		booking.PaymentPlan{},
		booking.PaymentDistribution{},
		hostel.HostelImage{},
		hostel.HostelFee{},
		hostel.HostelAgreementTemplate{},

		// maintenance
		maintenance.HostelMaintenanceRequest{},
		maintenance.HostelMaintenanceRequestImage{},
		maintenance.WorkOrder{},
		maintenance.WorkOrderImage{},
		maintenance.WorkOrderComment{},

		// notification
		notification.Notification{},

		// vendor
		vendor.Vendor{},

		// Review
		review.HostelReview{},
		review.HostelManagerReview{},
		review.VendorReview{},
	)
	if err != nil {
		fmt.Println(err)
		//log.Fatal().Err(err).Msg("failed to migrate")
	}
}
