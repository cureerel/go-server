package middleware

import (
	"net/http"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

// WriterContent allows primary writers and partners (marketplace) to manage posts — not reviewers-only accounts.
func WriterContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIface, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		r, _ := roleIface.(string)
		switch r {
		case entity.RoleWriter, entity.RolePartner, entity.RoleAdmin, entity.RoleSuperAdmin:
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "writer or partner role required"})
	}
}
