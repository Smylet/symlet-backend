package jobs

import (
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

// InitializeCronJobs initializes and starts all cron jobs in the application.
func InitializeCronJobs(cronScheduler *cron.Cron, db *gorm.DB, task worker.TaskDistributor) error {
	// Initialize the cron scheduler.

	// Register your cron jobs here.
	err := registerEmailVerificationReminderJob(cronScheduler, db, task)
	if err != nil {
		return err
	}
	// Add other cron jobs here as needed.

	// Start the cron scheduler.
	cronScheduler.Start()

	return nil
}
