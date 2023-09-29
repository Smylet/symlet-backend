package db

import (
	// "context"
	// "fmt"
	// "log"

	// "ariga.io/atlas-go-sdk/atlasexec"

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

func GetMigrateModels() []interface{} {
	return []interface{}{
		// List of models to be migrated
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
	}
}

// func Migrate(config utils.Config) {
// 	// Define the execution context, supplying a migration directory
// 	// and potentially an `atlas.hcl` configuration file using `atlasexec.WithHCL`.
// 	// Initialize the client.
// 	client, err := atlasexec.NewClient(fmt.Sprintf("%v/migrations", config.BasePath), "atlas")
// 	if err != nil {
// 		log.Fatalf("failed to initialize client: %v", err)
// 	}
	
// 	// Run `atlas migrate apply` on a SQLite database under /tmp.
// 	res, err := client.Apply(context.Background(), &atlasexec.ApplyParams{
// 		URL: fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
// 			config.DBUser,
// 			config.DBPass,
// 			config.DBHost,
// 			config.DBPort,
// 			config.DBName,
// 		),
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to apply migrations: %v", err)
// 	}
// 	fmt.Printf("Applied %d migrations\n", len(res.Applied))
// }
