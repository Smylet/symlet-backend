package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/mail/templates"
	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmailReminder = "task:send_verify_email_reminder"

type PayloadSendVerifyEmailReminder struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmailReminder(
	ctx context.Context,
	payload *PayloadSendVerifyEmailReminder,
	opts ...asynq.Option,
) error {
	logger := common.NewLogger()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmailReminder, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Printf(ctx, "enqueued task: %s", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmailReminder(ctx context.Context, task *asynq.Task) error {
	logger := common.NewLogger()

	var payload PayloadSendVerifyEmailReminder
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	data := mail.Data{
		Subject:       "Complete Your Registration: Email Verification Pending",
		To:            []string{payload.Email},
		From:          processor.config.EmailSenderName,
		Cc:            []string{},
		Bcc:           []string{},
		UserName:      payload.UserName,
		EmailTemplate: templates.EmailVerificationReminderTemplate,
		TemplateName:  templates.EmailVerificationReminder,
		Email:         payload.Email,
	}

	content, err := mail.GeneratePersonalizedEmail(data)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	data.Content = content

	errs := processor.mailer.SendEmail([]mail.Data{data})
	if len(errs) > 0 {
		logger.Error(errs[0].Error())
		return fmt.Errorf("failed to send email: %w", errs[0])
	}

	logger.Printf(ctx, "sent verification reminder to %s", payload.Email)
	return nil
}
