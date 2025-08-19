package handlers

import (
	"biletter-service/internal/middleware"
	"biletter-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) CreateBooking(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.services.Booking.CreateBooking(&req, currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to create booking", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *Handlers) ListBookings(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	bookings, err := h.services.Booking.GetBookingsByUser(currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to get bookings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *Handlers) CancelBooking(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Booking.CancelBooking(&req, currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to cancel booking", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) InitiatePayment(c *gin.Context) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentURL, err := h.services.Payment.InitiatePayment(&req, currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to initiate payment", zap.Error(err))
		c.JSON(http.StatusConflict, nil)
		return
	}

	c.Header("Location", paymentURL)
	c.JSON(http.StatusFound, nil)
}
