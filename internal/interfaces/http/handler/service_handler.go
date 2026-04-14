// internal/interfaces/http/handler/service_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	svc *service.ServiceService
}

func NewServiceHandler(svc *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

// GET /api/services — public

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




// POST /api/services/:id/live — 
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

// POST /api/services/:id/pause 
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

//  mapper 

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