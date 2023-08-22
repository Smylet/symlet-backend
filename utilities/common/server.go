package common

import (
	"log"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config) (*Server, error) {

	server := &Server{}

	r := gin.Default()
	users.RegisterRoutes(r)

	server.router = r
	return server, nil
}

// Start runs the HTTP server on a specific address.

func RunGinServer(config utils.Config) {
	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
