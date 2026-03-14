// internal/interfaces/http/middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			fmt.Println("[AuthMiddleware] missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			fmt.Println("[AuthMiddleware] invalid Authorization header format:", header)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header must be: Bearer <token>",
			})
			return
		}

		claims, err := authService.ValidateAccessToken(parts[1])
		if err != nil {
			fmt.Println("[AuthMiddleware] token validation failed:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		fmt.Printf("[AuthMiddleware] user authenticated: id=%s, email=%s, role=%s\n",
			claims.UserID, claims.Email, claims.Role)

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleMiddleware enforces a minimum role using the full hierarchy:

//	user(1) < writer(2) < partner(3) < worker(4) < admin(5) < superadmin(6)


func RoleMiddleware(minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "role not found — is AuthMiddleware applied first?",
			})
			return
		}

		userRole, ok := roleIface.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid role type in context",
			})
			return
		}

		u := &entity.User{Role: userRole}
		if !u.HasRole(minRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":     "insufficient permissions",
				"required":  minRole,
				"your_role": userRole,
			})
			return
		}

		c.Next()
	}
}

// SelfOrRoleMiddleware allows the request when EITHER condition is true:
//  1. The token's user_id matches the :id route param (you are the resource owner), OR
//  2. The token's role is >= minRole.
//
// Useful for routes like GET /users/:id where users can read their own profile
// but admins can read anyone's.
//
// Usage:
//
//	protected.GET("/users/:id",
//	    middleware.SelfOrRoleMiddleware(entity.RoleAdmin),
//	    userHandler.GetUserByID,
//	)
func SelfOrRoleMiddleware(minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenIDIface, _ := c.Get("user_id")
		tokenIDStr, _ := tokenIDIface.(string)

		paramID := c.Param("id")
		if paramID != "" && tokenIDStr == paramID {
			// Exact match — user is accessing their own resource.
			c.Next()
			return
		}

		// Fallback: check role hierarchy.
		roleIface, _ := c.Get("role")
		userRole, _ := roleIface.(string)
		u := &entity.User{Role: userRole}
		if !u.HasRole(minRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}


func OwnerOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenIDIface, _ := c.Get("user_id")
		tokenIDStr, _ := tokenIDIface.(string)

		ownerIface, exists := c.Get("resource_owner_id")
		if exists {
			ownerStr := fmt.Sprintf("%v", ownerIface)
			if tokenIDStr == ownerStr {
				c.Next()
				return
			}
		}

		roleIface, _ := c.Get("role")
		userRole, _ := roleIface.(string)
		u := &entity.User{Role: userRole}
		if !u.HasRole(entity.RoleAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you don't own this resource",
			})
			return
		}

		c.Next()
	}
}