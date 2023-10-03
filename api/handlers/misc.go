package handlers

import (
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) HealthCheck(c *gin.Context) {
	utils.RespondWithSuccess(c, 200, "OK", "OK")
}
