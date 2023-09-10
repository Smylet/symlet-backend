package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"

	"github.com/Smylet/symlet-backend/api/handlers"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/mail"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/Smylet/symlet-backend/utilities/worker"
)


// @title           Smylet API
// @version         1.0
// @description     This are the Smylet APP APIs.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization


// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if config.Environment == "development" {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	fmt.Println(config)
	database, redisOption := db.GetDB(config),
		asynq.RedisClientOpt{
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
