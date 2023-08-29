package db

import (
	"github.com/Smylet/symlet-backend/api/booking"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/maintenance"
	"github.com/Smylet/symlet-backend/api/notification"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/review"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/vendor"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		reference.ReferenceHostelAmmenities{},
		reference.ReferenceUniversity{},
		student.Student{},
		//hostel
		hostel.Hostel{},
		booking.HostelStudent{},
		booking.HostelBooking{},
		booking.CustomPaymentPlan{},
		booking.PaymentDistribution{},
		hostel.HostelImage{},
		hostel.HostelFee{},
		hostel.HostelAgreementTemplate{},

		//maintenance
		maintenance.HostelMaintenanceRequest{},
		maintenance.HostelMaintenanceRequestImage{},
		maintenance.WorkOrder{},
		maintenance.WorkOrderImage{},
		maintenance.WorkOrderComment{},

		//notification
		notification.Notification{},

		//vendor
		vendor.Vendor{},

		//Review
		review.HostelReview{},
		review.HostelManagerReview{},
		review.VendorReview{},
	)
	if err != nil {
		panic(err)
	}
}
