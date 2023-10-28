package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

func (server *Server) CreateStudent(c *gin.Context) {
	var studentSerializer student.StudentSerializer
	if err := c.ShouldBindJSON(&studentSerializer); err != nil {
		utils.RespondWithError(c, 400, err.Error(), "Invalid request payload")
		return
	}
	err := studentSerializer.Create(c, server.db, server.session)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "unable to create student")
		return
	}
	utils.RespondWithSuccess(c, 201, studentSerializer.Response(), "student created successfully")
}
