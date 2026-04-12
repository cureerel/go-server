package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/gin-gonic/gin"
)

type CoinHandler struct {
	coinSvc *service.CoinService
}

func NewCoinHandler(coinSvc *service.CoinService) *CoinHandler {
	return &CoinHandler{coinSvc: coinSvc}
}

func (h *CoinHandler) GetBalance(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	bal, err := h.coinSvc.Balance(c.Request.Context(), uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, gin.H{"balance": bal})
}

func (h *CoinHandler) PurchaseMembership(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.coinSvc.PurchaseMembership(c.Request.Context(), uid, req.Plan)
	if err != nil {
		respondErr(c, err)
		return
	}

	respond(c, gin.H{"message": "membership purchased successfully"})
}

