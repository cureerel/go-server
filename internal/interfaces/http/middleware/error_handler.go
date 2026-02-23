package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/cureerel/gotemplate/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware catches errors from handlers and returns proper JSON responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		// If no errors occurred, continue
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Check if the error is an AppError
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			// Log internal error if exists
			if appErr.Err != nil {
				log.Printf("AppError internal: %v", appErr.Err)
			}

			// Respond with safe public message
			appErr.RespondAndAbort(c)
			return
		}

		// Log unknown/unexpected errors
		log.Printf("Unexpected error: %v", err)

		// Default 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}