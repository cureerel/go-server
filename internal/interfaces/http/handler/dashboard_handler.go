// internal/interfaces/http/handler/dashboard_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GET /api/dashboard — role-dispatched single endpoint
func (h *DashboardHandler) Get(c *gin.Context) {
	uid, _ := getUID(c)
	role   := getRole(c)
	ctx    := c.Request.Context()

	switch {
	case hasRole(c, entity.RoleAdmin):
		data, err := h.svc.AdminDashboard(ctx)
		if err != nil { respondErr(c, err); return }
		c.JSON(http.StatusOK, gin.H{"role": role, "data": data})

	case role == entity.RoleWorker:
		data, err := h.svc.WorkerDashboard(ctx, uid)
		if err != nil { respondErr(c, err); return }
		c.JSON(http.StatusOK, gin.H{"role": role, "data": data})

	case role == entity.RolePartner:
		data, err := h.svc.PartnerDashboard(ctx, uid)
		if err != nil { respondErr(c, err); return }
		c.JSON(http.StatusOK, gin.H{"role": role, "data": data})

	case hasRole(c, entity.RoleWriter):
		data, err := h.svc.WriterDashboard(ctx, uid)
		if err != nil { respondErr(c, err); return }
		c.JSON(http.StatusOK, gin.H{"role": role, "data": data})

	default:
		data, err := h.svc.UserDashboard(ctx, uid)
		if err != nil { respondErr(c, err); return }
		c.JSON(http.StatusOK, gin.H{"role": role, "data": data})
	}
}

// GET /api/dashboard/user    — explicit per-role endpoints (useful for frontend)
func (h *DashboardHandler) UserView(c *gin.Context) {
	uid, _ := getUID(c)
	data, err := h.svc.UserDashboard(c.Request.Context(), uid)
	if err != nil { respondErr(c, err); return }
	respond(c, data)
}

func (h *DashboardHandler) WriterView(c *gin.Context) {
	uid, _ := getUID(c)
	data, err := h.svc.WriterDashboard(c.Request.Context(), uid)
	if err != nil { respondErr(c, err); return }
	respond(c, data)
}

func (h *DashboardHandler) PartnerView(c *gin.Context) {
	uid, _ := getUID(c)
	data, err := h.svc.PartnerDashboard(c.Request.Context(), uid)
	if err != nil { respondErr(c, err); return }
	respond(c, data)
}

func (h *DashboardHandler) WorkerView(c *gin.Context) {
	uid, _ := getUID(c)
	data, err := h.svc.WorkerDashboard(c.Request.Context(), uid)
	if err != nil { respondErr(c, err); return }
	respond(c, data)
}

func (h *DashboardHandler) AdminView(c *gin.Context) {
	data, err := h.svc.AdminDashboard(c.Request.Context())
	if err != nil { respondErr(c, err); return }
	respond(c, data)
}