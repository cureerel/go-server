package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

// ----------------- Requests -----------------

type activateMembershipRequest struct {
	Plan entity.MembershipPlan `json:"plan" binding:"required"`
}

type upgradeMembershipRequest struct {
	Plan entity.MembershipPlan `json:"plan" binding:"required"`
}

// POST /memberships/activate
func (h *MembershipHandler) Activate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req activateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

	membership, err := h.membershipService.Activate(c.Request.Context(), uint(uid), req.Plan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toMembershipResponse(membership)})
}

// GET /memberships/me
func (h *MembershipHandler) GetMine(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

	membership, err := h.membershipService.GetByUserID(c.Request.Context(), uint(uid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toMembershipResponse(membership)})
}

// POST /memberships/upgrade
func (h *MembershipHandler) Upgrade(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req upgradeMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

	membership, err := h.membershipService.Upgrade(c.Request.Context(), uint(uid), req.Plan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toMembershipResponse(membership)})
}

// DELETE /memberships/cancel
func (h *MembershipHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

	if err := h.membershipService.Cancel(c.Request.Context(), uint(uid)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membership cancelled"})
}

// Helper

func toMembershipResponse(m *entity.Membership) dto.MembershipResponse {
	return dto.MembershipResponse{
		ID:        m.ID,
		UserID:    m.UserID,
		Plan:      string(m.Plan),
		Status:    string(m.Status),
		StartsAt:  m.StartsAt.Format("2006-01-02T15:04:05Z"),
		ExpiresAt: m.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
