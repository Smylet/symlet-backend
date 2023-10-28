package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/mail/templates"
	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserName            string `json:"username"`
	UserID              uint   `json:"user_id"`
	SecretCode          uint   `json:"secret_code"`
	VerificationEmailID uint   `json:"ver_email_id"`
	Email               string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	logger := common.NewLogger()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Printf(ctx, "enqueued task: %s", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	logger := common.NewLogger()

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// Convert payload.SecretCode to a string
	secretCodeStr := strconv.Itoa(int(payload.SecretCode))

	// TODO: replace this URL with an environment variable that points to a front-end page - GET /users/confirm-email
	verifyUrl := fmt.Sprintf("http://%s/users/confirm-email?user_id=%d&ver_email_id=%d&secret_code=%s",
		processor.config.HTTPServerAddress,
		payload.UserID,
		payload.VerificationEmailID,
		secretCodeStr)

	data := mail.Data{
		Subject:       "Welcome to Smylet!",
		To:            []string{payload.Email},
		From:          processor.config.EmailSenderName,
		Cc:            []string{},
		Bcc:           []string{},
		UserName:      payload.UserName,
		Url:           verifyUrl,
		EmailTemplate: templates.RegistrationTemplate,
		TemplateName:  templates.Registration,
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

	logger.Printf(ctx, "sent verification email to %s", payload.Email)
	return nil
}
