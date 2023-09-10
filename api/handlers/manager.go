package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/utilities/utils"
)


// @Summary Create a hostel manager
// @Description Create a new hostel manager
// @Tags Hostel Managers
// @Accept json
// @Produce json
// @Param hostel_manager body manager.HostelManagerSerializer true "Hostel manager object to create"
// @Success 201 {object} manager.HostelManagerSerializer
// @Failure 400 {object} utils.ErrorMessage
// @Failure 500 {object} utils.ErrorMessage
// @Router /hostel-managers [post]
func (server *Server) CreateHostelManager(c *gin.Context) {
	var HostelManagerSerializer manager.HostelManagerSerializer

	if err := c.ShouldBind(&HostelManagerSerializer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := HostelManagerSerializer.CreateTx(c, server.db)

	if err != nil {
		utils.RespondWithError(c, 500, err.Error(), "Failed to create hostel manager")
		return
	}

	utils.RespondWithSuccess(c, 201, HostelManagerSerializer.Response(), "Hostel manager created successfully")
}
