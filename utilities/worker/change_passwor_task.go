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

const TaskSendChangePasswordEmail = "task:send_change_password_email"

type PayloadSendChangePasswordEmail struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendChangePasswordEmail(
	ctx context.Context,
	payload *PayloadSendChangePasswordEmail,
	opts ...asynq.Option,
) error {
	logger := common.NewLogger()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendChangePasswordEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Printf(ctx, "enqueued task: %s", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendChangePasswordEmail(ctx context.Context, task *asynq.Task) error {
	logger := common.NewLogger()

	var payload PayloadSendChangePasswordEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	data := mail.Data{
		Subject:       "Change your password",
		From:          processor.config.EmailSenderName,
		To:            []string{payload.Email},
		Email:         payload.Email,
		UserName:      payload.UserName,
		EmailTemplate: templates.ChangePasswordTemplate,
		TemplateName:  templates.ChangePassword,
		Cc:            []string{},
		Bcc:           []string{},
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

	logger.Printf(ctx, "sent email to %s", payload.Email)

	return nil
}
