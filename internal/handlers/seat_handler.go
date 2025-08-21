package handlers

import (
	"biletter-service/internal/middleware"
	"biletter-service/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) ListSeats(c *gin.Context) {
	eventIDStr := c.Query("event_id")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
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

	seats, err := h.services.Seat.GetSeatsByEvent(eventID, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get seats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
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
		c.JSON(http.StatusInsufficientStorage, nil)
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
		c.JSON(http.StatusInsufficientStorage, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) FillSeats(c *gin.Context) {
	h.services.Seat.FillSeats()

	c.JSON(http.StatusOK, nil)
}
