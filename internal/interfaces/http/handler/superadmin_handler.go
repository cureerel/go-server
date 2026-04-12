// internal/interfaces/http/handler/superadmin_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type SuperAdminHandler struct {
	svc *service.SuperAdminService
}

func NewSuperAdminHandler(svc *service.SuperAdminService) *SuperAdminHandler {
	return &SuperAdminHandler{svc: svc}
}

// GET /api/superadmin/stats
func (h *SuperAdminHandler) Stats(c *gin.Context) {
	stats, err := h.svc.PlatformStats(c.Request.Context())
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, stats)
}

// PATCH /api/superadmin/users/:id/role
func (h *SuperAdminHandler) SetRole(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.SetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	callerID, _ := getUID(c)
	if err := h.svc.SetRole(c.Request.Context(), userID, req.Role, callerID); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

// GET /api/superadmin/upgrades
func (h *SuperAdminHandler) GetPendingUpgrades(c *gin.Context) {
	page, limit := paginate(c)
	reqs, total, err := h.svc.GetPendingUpgrades(c.Request.Context(), page, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.UpgradeRequestResponse, len(reqs))
	for i := range reqs {
		list[i] = toUpgradeResponse(&reqs[i])
	}
	c.JSON(http.StatusOK, dto.UpgradeRequestListResponse{
		Data: list, Total: total, Page: page, Limit: limit,
	})
}

// POST /api/superadmin/upgrades/:id/review
func (h *SuperAdminHandler) ReviewUpgrade(c *gin.Context) {
	reqID, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.ReviewUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reviewerID, _ := getUID(c)
	if err := h.svc.ReviewUpgrade(c.Request.Context(), reqID, req.Approve, reviewerID); err != nil {
		respondErr(c, err)
		return
	}
	msg := "upgrade request rejected"
	if req.Approve {
		msg = "upgrade request approved"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// ── Self-service upgrade request (any auth user) ──────────────

// POST /api/upgrade-requests
func (h *SuperAdminHandler) RequestUpgrade(c *gin.Context) {
	var req dto.RequestUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	result, err := h.svc.RequestUpgrade(c.Request.Context(), uid, req.ToRole)
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toUpgradeResponse(result))
}

// GET /api/upgrade-requests/me
func (h *SuperAdminHandler) GetMyUpgradeRequest(c *gin.Context) {
	uid, _ := getUID(c)
	req, err := h.svc.GetMyUpgradeRequest(c.Request.Context(), uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toUpgradeResponse(req))
}

// ── mapper ────────────────────────────────────────────────────

func toUpgradeResponse(r *entity.UpgradeRequest) dto.UpgradeRequestResponse {
	resp := dto.UpgradeRequestResponse{
		ID:         r.ID,
		UserID:     r.UserID,
		FromRole:   r.FromRole,
		ToRole:     r.ToRole,
		Status:     r.Status,
		ReviewedBy: r.ReviewedBy,
		CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if r.ReviewedAt != nil {
		s := r.ReviewedAt.Format("2006-01-02T15:04:05Z")
		resp.ReviewedAt = &s
	}
	return resp
}