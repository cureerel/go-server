package handler

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/cureerel/gotemplate/internal/application/service"
    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/gin-gonic/gin"
)

type MembershipHandler struct {
    membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
    return &MembershipHandler{membershipService: membershipService}
}

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
    c.JSON(http.StatusCreated, gin.H{"data": membership})
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
    c.JSON(http.StatusOK, gin.H{"data": membership})
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
    c.JSON(http.StatusOK, gin.H{"data": membership})
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