package jobs

import (
	"time"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/hibiken/asynq"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

const chunkSize = 100 // Number of users processed in a single batch

func registerEmailVerificationReminderJob(scheduler *cron.Cron, db *gorm.DB, task worker.TaskDistributor) error {
	logger := common.NewLogger()

	// Add a cron job to send the email reminder after 24 hours.
	err := scheduler.AddFunc("@hourly", func() {

		var users []users.User
		if err := db.Where("is_email_confirmed = ? AND email_reminder_count < ?", false, 3).Find(&users).Error; err != nil {
			logger.Error("Failed to query users:", err)
			return
		}

		// Process users in chunks
		for i := 0; i < len(users); i += chunkSize {
			end := i + chunkSize
			if end > len(users) {
				end = len(users)
			}
			go processUserChunk(users[i:end], db, task)
		}
	})

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func processUserChunk(users []users.User, db *gorm.DB, task worker.TaskDistributor) {
	logger := common.NewLogger()

	for _, user := range users {
		tx := db.Begin()

		payload := worker.PayloadSendVerifyEmailReminder{
			UserName: user.UserName,
			Email:    user.Email,
		}

		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(10 * time.Second),
			asynq.Queue(worker.QueueCritical),
		}

		if err := task.DistributeTaskSendVerifyEmailReminder(context.Background(), &payload, opts...); err != nil {
			logger.Error("Failed to distribute task:", err)
			tx.Rollback()
			continue
		}

		user.EmailReminderCount++
		if err := tx.Save(&user).Error; err != nil {
			logger.Error("Failed to update EmailReminderCount for user:", user.ID, "Error:", err)
			tx.Rollback()
			continue
		}

		tx.Commit()
	}
}
