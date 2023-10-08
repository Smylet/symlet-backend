package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Smylet/symlet-backend/api/vendor"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

func (server *Server) CreateVendor(c *gin.Context) {
	var vendorSerializer vendor.VendorSerializer
	if err := common.CustomBinder(c, vendorSerializer); err != nil {
		utils.RespondWithError(
			c, 400, err.Error(), fmt.Sprint("unable to bind request body to serializer"),
		)
		return
	}
	if err := vendorSerializer.Create(c, server.db, server.session); err != nil {
		utils.RespondWithError(c, 400, err.Error(), fmt.Sprint("unable to create vendor"))
		return
	}
	utils.RespondWithSuccess(c, 201, vendorSerializer.Response(),"vendor created successfully")
}

func (server *Server) GetVendor(c *gin.Context) {
	var vendorSerializer vendor.VendorSerializer
	uidString := c.Param("uid")
	if uidString == "" {
		utils.RespondWithError(
			c, 400, "uid is required", "uid is required",
		)
	}
	if err := vendorSerializer.Get(c, server.db, uidString); err != nil {
		utils.RespondWithError(c, 400, err.Error(), fmt.Sprint("unable to get vendor"))
	}
	utils.RespondWithSuccess(c, 200, vendorSerializer.Response(),"vendor retrieved successfully")
}
