package handlers

import (
	"biletter-service/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) ListSeats(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	seats, err := h.services.Seat.GetSeatsByEvent(eventID)
	if err != nil {
		h.logger.Error("Failed to get seats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func (h *Handlers) SelectSeat(c *gin.Context) {
	var req models.SelectSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Seat.SelectSeat(&req)
	if err != nil {
		h.logger.Error("Failed to select seat", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat selected successfully"})
}

func (h *Handlers) ReleaseSeat(c *gin.Context) {
	var req models.ReleaseSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Seat.ReleaseSeat(&req)
	if err != nil {
		h.logger.Error("Failed to release seat", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat released successfully"})
}