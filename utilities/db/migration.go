package db

import (
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
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(

		users.User{},
		users.Profile{},
		users.VerificationEmail{},
		users.Session{},
		student.Student{},

		manager.HostelManager{},
		reference.ReferenceHostelAmmenities{},
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
		log.Fatal().Err(err).Msg("failed to migrate")
	}
}
