package handlers

import (
	"biletter-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) InitiatePayment(c *gin.Context) {
	var req models.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Payment.InitiatePayment(&req)
	if err != nil {
		h.logger.Error("Failed to initiate payment", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment initiated successfully"})
}