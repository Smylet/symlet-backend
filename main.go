package main

import (
	"log"
	"os"

	"github.com/Smylet/symlet-backend/api/handlers"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
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

	server, err := handlers.NewServer(config, database)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
