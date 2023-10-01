package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Smylet/symlet-backend/api/handlers"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if config.Environment == "development" {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Define a channel to communicate errors from goroutines.
	errCh := make(chan error, 2)

	var database *gorm.DB
	go func() {
		db, err := db.GetDB(config)
		if err != nil {
			errCh <- fmt.Errorf("failed to connect to database: %w", err)
			if db != nil {
				db.Close()
			}
			return
		}

		if db == nil {
			errCh <- errors.New("failed to connect to database - db is nil")
			return
		}

		database = db.GormDB()
		errCh <- nil // Send nil to signify successful completion.
	}()

	redisOption := asynq.RedisClientOpt{Addr: config.RedisAddress}
	var awsSession *session.Session
	go func() {
		awsSession, err = common.CreateAWSSession(&config)
		if err != nil {
			errCh <- fmt.Errorf("failed to create AWS session: %w", err)
			return
		}
		errCh <- nil // Send nil to signify successful completion.
	}()

	// Wait for both goroutines to complete and check for errors.
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			logger.Fatal().Err(err)
			os.Exit(1)
		}
	}

	mailer, task := mail.NewSESEmailSender(
		config.EmailSenderAddress,
		awsSession,
		config),
		worker.NewRedisTaskDistributor(redisOption)

	go worker.RunTaskProcessor(config, redisOption, database, mailer)

	server, err := handlers.NewServer(config, database, task, mailer, awsSession)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}

	mail.Cleanup()
}
