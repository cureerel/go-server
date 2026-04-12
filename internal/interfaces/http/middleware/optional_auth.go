package middleware

import (
	"strings"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/gin-gonic/gin"
)

// OptionalAuth sets user_id / email / role when a valid Bearer token is present; never aborts.
func OptionalAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}
		claims, err := authService.ValidateAccessToken(parts[1])
		if err != nil {
			c.Next()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
