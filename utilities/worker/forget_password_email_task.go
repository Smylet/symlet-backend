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

const TaskSendForgetPasswordEmail = "task:send_forget_password_email"

type PayloadSendForgetPasswordEmail struct {
	UserName   string `json:"username"`
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendForgetPasswordEmail(
	ctx context.Context,
	payload *PayloadSendForgetPasswordEmail,
	opts ...asynq.Option,
) error {
	logger := common.NewLogger()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendForgetPasswordEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Printf(ctx, "enqueued task: %s", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendForgetPasswordEmail(ctx context.Context, task *asynq.Task) error {
	logger := common.NewLogger()

	var payload PayloadSendForgetPasswordEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	resetLink := fmt.Sprintf("http://%s/reset-password?reset_token=%s",
		processor.config.HTTPServerAddress,
		payload.ResetToken)

	data := mail.Data{
		Subject:       "Reset your password",
		To:            []string{payload.Email},
		From:          processor.config.EmailSenderName,
		Email:         payload.Email,
		UserName:      payload.UserName,
		Url:           resetLink,
		EmailTemplate: templates.ResetPasswordTemplate,
		TemplateName:  templates.ResetPassword,
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
