package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAnyRole allows the request when the JWT role matches one of allowed (exact match).
func RequireAnyRole(allowed ...string) gin.HandlerFunc {
	set := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		set[a] = struct{}{}
	}
	return func(c *gin.Context) {
		roleIface, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		r, _ := roleIface.(string)
		if _, ok := set[r]; ok {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
