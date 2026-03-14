// internal/interfaces/http/handler/service_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	svc *service.ServiceService
}

func NewServiceHandler(svc *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

// GET /api/services — public
// ?status=live (default) | pending | approved | paused
// ?search=keyword
func (h *ServiceHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.DefaultQuery("status", entity.ServiceStatusLive)
	search := c.Query("search")

	svcs, total, err := h.svc.GetAll(c.Request.Context(), page, limit, status, search)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.ServiceResponse, len(svcs))
	for i := range svcs {
		list[i] = toServiceResponse(&svcs[i])
	}
	c.JSON(http.StatusOK, dto.ServiceListResponse{
		Data: list, Total: total, Page: page, Limit: limit,
	})
}

// GET /api/services/:id — public
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	svc, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	go func() { _ = h.svc.RecordView(c.Request.Context(), id, c.ClientIP(), c.Request.UserAgent()) }()
	respond(c, toServiceResponse(svc))
}

// GET /api/services/mine — partner+ (own services, any status)
func (h *ServiceHandler) GetMine(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	page, limit := paginate(c)
	svcs, total, err := h.svc.GetByOwner(c.Request.Context(), uid, page, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.ServiceResponse, len(svcs))
	for i := range svcs {
		list[i] = toServiceResponse(&svcs[i])
	}
	c.JSON(http.StatusOK, dto.ServiceListResponse{
		Data: list, Total: total, Page: page, Limit: limit,
	})
}

// POST /api/services — partner+
func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	svc, err := h.svc.Create(c.Request.Context(), service.CreateServiceInput{
		OwnerID:       uid,
		Title:         req.Title,
		Description:   req.Description,
		PriceUSDCents: req.PriceUSDCents,
		CoverImageURL: req.CoverImageURL,
		CoverImageKey: req.CoverImageKey,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toServiceResponse(svc))
}

// PUT /api/services/:id — partner+ (service enforces ownership)
func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	svc, err := h.svc.Update(c.Request.Context(), id, uid, service.UpdateServiceInput{
		Title:         req.Title,
		Description:   req.Description,
		PriceUSDCents: req.PriceUSDCents,
		CoverImageURL: req.CoverImageURL,
		CoverImageKey: req.CoverImageKey,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toServiceResponse(svc))
}

// DELETE /api/services/:id — partner+ (service enforces ownership or admin)
func (h *ServiceHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if err := h.svc.Delete(c.Request.Context(), id, uid, getRole(c)); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service deleted"})
}

// POST /api/services/:id/approve — admin+
func (h *ServiceHandler) Approve(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.svc.Approve(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service approved"})
}

// POST /api/services/:id/reject — admin+
func (h *ServiceHandler) Reject(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.svc.Reject(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service rejected"})
}

// POST /api/services/:id/live — partner+ (own, must be approved first)
func (h *ServiceHandler) SetLive(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if err := h.svc.SetLive(c.Request.Context(), id, uid); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service is now live"})
}

// POST /api/services/:id/pause — partner+ (own) or admin
func (h *ServiceHandler) Pause(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if err := h.svc.Pause(c.Request.Context(), id, uid, getRole(c)); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service paused"})
}

// ── mapper ────────────────────────────────────────────────────

func toServiceResponse(s *entity.Service) dto.ServiceResponse {
	return dto.ServiceResponse{
		ID:            s.ID,
		OwnerID:       s.OwnerID,
		Title:         s.Title,
		Description:   s.Description,
		PriceUSDCents: s.PriceUSDCents,
		PriceUSD:      s.PriceUSD(),
		Status:        s.Status,
		CoverImageURL: s.CoverImageURL,
		ViewsTotal:    s.ViewsTotal,
		CreatedAt:     s.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}