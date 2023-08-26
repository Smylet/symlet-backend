package worker

import (
	"context"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	// ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db     *gorm.DB
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db *gorm.DB, mailer mail.EmailSender) TaskProcessor {
	logger := common.NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
		db:     db,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	// mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
