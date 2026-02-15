package apperror

import (
    "github.com/gin-gonic/gin"
)

// Respond writes error response to gin context
func (e *AppError) Respond(c *gin.Context) {
    c.JSON(e.HTTPStatus(), gin.H{
        "error": e.PublicError(),
        "code":  e.Code,
    })
}

// RespondAndAbort stops further handlers
func (e *AppError) RespondAndAbort(c *gin.Context) {
    e.Respond(c)
    c.Abort()
}