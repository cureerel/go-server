// pkg/apperror/gin_handler.go
package apperror

import "github.com/gin-gonic/gin"

// Respond writes a JSON error response using AppError if available,
// otherwise falls back to a generic 500. Used by all handlers via respondErr().
func Respond(c *gin.Context, err error) {
	c.JSON(HTTPStatus(err), gin.H{"error": PublicMessage(err)})
}
