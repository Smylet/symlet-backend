package main

import (
	"log"
	"os"

	"github.com/Smylet/symlet-backend/api/handlers"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if config.Environment == "development" {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	db, err := db.GetDB(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		os.Exit(1)
	}

	if db == nil {
		logger.Fatal().Msg("failed to connect to database - db is nil")
		os.Exit(1)
	}
	database := db.GormDB()

	redisOption := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	awsSession, err := common.CreateAWSSession(&config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create AWS session")
	}

	mailer, task := mail.NewSESEmailSender(
		config.EmailSenderAddress,
		awsSession,
		config),
		worker.NewRedisTaskDistributor(redisOption)

	go worker.RunTaskProcessor(config, redisOption,
		database, mailer)

	server, err := handlers.NewServer(config, database,
		task, mailer, awsSession)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}

	mail.Cleanup()
}
