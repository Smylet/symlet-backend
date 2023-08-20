package main

import (
	"log"

	"github.com/gin-gonic/gin"
	users "github.com/r-scheele/project/api/user"
)

func main() {
	r := gin.Default()

	// Register routes
	users.RegisterRoutes(r)

	log.Println("Starting server on port 8080")
	// Start the Gin server
	r.Run() // default to :8080
}
