package main

import (
	"log"

	users "github.com/The-CuriousX/project/api/user"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register routes
	users.RegisterRoutes(r)

	log.Println("Starting server on port 8080")
	// Start the Gin server
	r.Run() // default to :8080
}
