// internal/interfaces/http/middleware/error_handler.go
package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.Err != nil {
				log.Printf("AppError internal: %v", appErr.Err)
			}
			c.JSON(appErr.HTTPStatus(), gin.H{"error": appErr.PublicError()})
			return
		}

		log.Printf("Unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
