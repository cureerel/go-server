// internal/interfaces/http/handler/ticket_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

// POST /api/tickets —  authenticated user
func (h *TicketHandler) Create(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	ticket, err := h.svc.Create(c.Request.Context(), service.CreateTicketInput{
		UserID:      uid,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toTicketResponse(ticket))
}

// GET /api/tickets/me — own tickets
func (h *TicketHandler) GetMine(c *gin.Context) {
	uid, _ := getUID(c)
	page, limit := paginate(c)
	status := c.Query("status")
	tickets, total, err := h.svc.GetMine(c.Request.Context(), uid, page, limit, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.TicketResponse, len(tickets))
	for i := range tickets {
		list[i] = toTicketResponse(&tickets[i])
	}
	c.JSON(http.StatusOK, dto.TicketListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// GET /api/tickets —
func (h *TicketHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.Query("status")
	priority := c.Query("priority")
	tickets, total, err := h.svc.GetAll(c.Request.Context(), page, limit, status, priority)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.TicketResponse, len(tickets))
	for i := range tickets {
		list[i] = toTicketResponse(&tickets[i])
	}
	c.JSON(http.StatusOK, dto.TicketListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// POST /api/tickets/:id/resolve
func (h *TicketHandler) Resolve(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.svc.Resolve(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket resolved"})
}

// POST /api/tickets/:id/close — owner or admin
func (h *TicketHandler) Close(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if err := h.svc.Close(c.Request.Context(), id, uid, getRole(c)); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket closed"})
}

// POST /api/tickets/:id/messages — owner or admin
func (h *TicketHandler) SendMessage(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	msg, err := h.svc.SendMessage(c.Request.Context(), id, uid, req.Message)
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toMessageResponse(msg))
}

// GET /api/tickets/:id/messages — owner or admin
func (h *TicketHandler) GetMessages(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	msgs, err := h.svc.GetMessages(c.Request.Context(), id, uid, getRole(c))
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.TicketMessageResponse, len(msgs))
	for i := range msgs {
		list[i] = toMessageResponse(&msgs[i])
	}
	respond(c, list)
}

// mappers

func toTicketResponse(t *entity.Ticket) dto.TicketResponse {
	r := dto.TicketResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Subject:     t.Subject,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		AssignedTo:  t.AssignedTo,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if t.ClosedAt != nil {
		s := t.ClosedAt.Format("2006-01-02T15:04:05Z")
		r.ClosedAt = &s
	}
	return r
}

func toMessageResponse(m *entity.TicketMessage) dto.TicketMessageResponse {
	return dto.TicketMessageResponse{
		ID:        m.ID,
		TicketID:  m.TicketID,
		SenderID:  m.SenderID,
		Message:   m.Message,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
