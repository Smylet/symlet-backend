package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

const (

	TaskSendVerifyEmail = "task:send_verify_email"
)


type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	arg := users.FindUserParams{
		User:           users.User{Username: payload.Username},
		IncludeProfile: false,
	}

	userRepo := users.NewUserRepository(processor.db)

	user, err := userRepo.FindUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := userRepo.CreateVerifyEmail(ctx, users.CreateVerifyEmailParams{
		UserID:     user.ID,
		Email:      user.Email,
		SecretCode: utils.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// TODO: replace this URL with an environment variable that points to a front-end page - GET /users/confirm-email
	verifyUrl := fmt.Sprintf("http://%s/users/confirm-email?user_id=%d&ver_email_id=%d&secret_code=%s",
		processor.config.HTTPServerAddress,
		verifyEmail.UserID,
		verifyEmail.ID, verifyEmail.SecretCode)

	data := mail.PersonalizedData{
		User:    user,
		Subject: "Welcome to Smylet!",
		Cc:      []string{},
		Bcc:     []string{},
		Others: map[string]interface{}{
			"VerificationLink": verifyUrl,
		},
	}

	content, err := mail.GeneratePersonalizedEmail(data)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	data.Content = content

	errs := processor.mailer.SendEmail([]mail.PersonalizedData{data})
	if len(errs) > 0 {
		return fmt.Errorf("failed to send email: %w", errs[0])
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}

