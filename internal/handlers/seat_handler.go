package handlers

import (
	"biletter-service/internal/middleware"
	"biletter-service/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) ListSeats(c *gin.Context) {
	eventIDStr := c.Query("event_id")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	rowStr := c.Query("row")
	status := c.Query("status")
	if eventIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}
	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSizeStr = "20"
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
	}
	if status != "" && !isValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status code"})
	}
	row, err := strconv.ParseInt(rowStr, 10, 64)

	seats, err := h.services.Seat.GetSeatsByEvent(eventID, status, row, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get seats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func isValidStatus(status string) bool {
	switch models.SeatStatus(status) {
	case models.SeatStatusFree, models.SeatStatusReserved, models.SeatStatusSold:
		return true
	}
	return false
}

func (h *Handlers) SelectSeat(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SelectSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Booking.SelectSeat(req.BookingID, req.SeatID, currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to select seat", zap.Error(err))
		if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInsufficientStorage, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) ReleaseSeat(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.ReleaseSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Booking.ReleaseSeat(req.SeatID, currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to release seat", zap.Error(err))
		if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInsufficientStorage, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) FillSeats(c *gin.Context) {
	h.services.Seat.FillSeats()

	c.JSON(http.StatusOK, nil)
}
