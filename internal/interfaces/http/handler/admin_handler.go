// internal/interfaces/http/handler/superadmin_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	svc *service.SuperAdminService
}

func NewSuperAdminHandler(svc *service.SuperAdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// GET /api/admin/stats
func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := h.svc.PlatformStats(c.Request.Context())
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, stats)
}
