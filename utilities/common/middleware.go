package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DatabaseTransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Begin a new transaction
		tx := db.Begin()
		if tx.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start a transaction",
			})
			return
		}

		// Set the transaction in the context so it can be accessed in request handlers
		ctx.Set("db_transaction", tx)

		// Continue processing the request
		ctx.Next()

		// Check for errors in the response status code
		if ctx.Writer.Status() >= http.StatusInternalServerError {
			// An HTTP 500 error occurred, rollback the database transaction
			if err := tx.Rollback().Error; err != nil {
				// Handle rollback error
				ctx.Error(err)
				return
			}
		} else {
			// Commit the transaction if no errors occurred
			if err := tx.Commit().Error; err != nil {
				// Handle commit error
				ctx.Error(err)
				return
			}
		}
	}
}
