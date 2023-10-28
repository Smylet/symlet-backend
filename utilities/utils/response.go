package utils

import (
	"github.com/gin-gonic/gin"
)

// Message represents a generic API response message.
type ErrorMessage struct {
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

type SuccessMessage struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// RespondWithError sends an error response with a given status code, error message, and a general message.
func RespondWithError(c *gin.Context, status int, errorMsg string, generalMsg string) {

	c.AbortWithStatusJSON(status, ErrorMessage{
		Error: errorMsg,
		Msg:   generalMsg,
	})
}

// RespondWithSuccess sends a success response with a given status code and data.
func RespondWithSuccess(c *gin.Context, status int, data interface{}, generalMsg string) {
	c.JSON(status, SuccessMessage{
		Data: data,
		Msg:  generalMsg,
	})
}
