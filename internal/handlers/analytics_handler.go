package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) GetAnalytics(c *gin.Context) {
	h.logger.Info("Get analytics endpoint called")

	// Получаем eventID из query параметра
	eventIDStr := c.Query("id")
	if eventIDStr == "" {
		h.logger.Error("Missing event ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing event ID parameter"})
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid event ID parameter", zap.String("id", eventIDStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID parameter"})
		return
	}

	analytics, err := h.services.Analytics.GetAnalytics(eventID)
	if err != nil {
		h.logger.Error("Failed to get analytics", zap.Int64("event_id", eventID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	h.logger.Info("Analytics retrieved successfully", zap.Int64("event_id", eventID))
	c.JSON(http.StatusOK, analytics)
}
