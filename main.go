package main

import (
	"log"
	"os"

	"github.com/Smylet/symlet-backend/api/handlers"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	if config.Environment == "development" {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	database := db.GetDB(config)

	mailer, task := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword),
		worker.NewRedisTaskDistributor(asynq.RedisClientOpt{
			Addr: config.RedisAddress,
		})

	server, err := handlers.NewServer(config, database, task, mailer)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
