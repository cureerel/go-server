// internal/interfaces/http/handler/helpers.go
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func respond(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func respondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func respondErr(c *gin.Context, err error) {
	c.JSON(apperror.HTTPStatus(err), gin.H{"error": apperror.PublicMessage(err)})
}

func getUID(c *gin.Context) (uint, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
	return uint(id), err == nil
}

func getRole(c *gin.Context) string {
	v, _ := c.Get("role")
	r, _ := v.(string)
	return r
}

func hasRole(c *gin.Context, min string) bool {
	u := &entity.User{Role: getRole(c)}
	return u.HasRole(min)
}

func parseID(c *gin.Context, param string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		return 0, apperror.NewBadRequest("invalid " + param)
	}
	return uint(id), nil
}

func paginate(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}
