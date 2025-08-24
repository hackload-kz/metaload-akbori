package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ResetData сбрасывает все данные броней и мест
func (h *Handlers) ResetData(c *gin.Context) {
	h.logger.Info("Reset data endpoint called")

	err := h.services.Reset.ResetAllData()
	if err != nil {
		h.logger.Error("Failed to reset data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset data"})
		return
	}

	h.logger.Info("Data reset completed successfully")
	c.JSON(http.StatusOK, gin.H{"message": "All data reset successfully"})
}
