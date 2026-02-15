package middleware

import (
    "errors"
    "github.com/cureerel/gotemplate/pkg/apperror"
    "github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Process request
        
        // Check if there are any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            // Check if it's our AppError
            var appErr *apperror.AppError
            if errors.As(err, &appErr) {
                c.JSON(appErr.HTTPStatus(), gin.H{
                    "error": appErr.PublicError(),
                })
                return
            }
            
            // Default error
            c.JSON(500, gin.H{"error": "internal server error"})
        }
    }
}