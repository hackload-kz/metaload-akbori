package handlers

import (
	"biletter-service/internal/models"
	"biletter-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	paymentService services.PaymentService
	logger         *zap.Logger
	validator      *validator.Validate
}

func NewPaymentHandler(paymentService services.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger,
		validator:      validator.New(),
	}
}

// PaymentNotifications handles webhook notifications from payment gateway
// POST /api/payments/notifications
func (h *PaymentHandler) PaymentNotifications(c *gin.Context) {
	var payload models.PaymentNotificationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Failed to bind payment notification payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&payload); err != nil {
		h.logger.Error("Payment notification payload validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	h.logger.Info("Payment webhook notification received",
		zap.String("paymentId", payload.PaymentID),
		zap.String("status", payload.Status),
		zap.Any("data", payload.Data))

	err := h.paymentService.ProcessPaymentNotification(&payload)
	if err != nil {
		h.logger.Error("Error processing payment webhook",
			zap.String("paymentId", payload.PaymentID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing webhook"})
		return
	}

	h.logger.Info("Successfully processed payment webhook",
		zap.String("paymentId", payload.PaymentID))

	c.String(http.StatusOK, "OK")
}

// PaymentSuccess handles redirect after successful payment
// GET /api/payments/success?orderId=123
func (h *PaymentHandler) PaymentSuccess(c *gin.Context) {
	orderId := c.Query("orderId")
	//if orderId == "" {
	//	h.logger.Error("Missing orderId parameter in payment success redirect")
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing orderId parameter"})
	//	return
	//}
	//
	//h.logger.Info("Payment success redirect received", zap.String("orderId", orderId))
	//
	//err := h.paymentService.NotifyPaymentSuccess(orderId)
	//if err != nil {
	//	h.logger.Error("Error processing payment success",
	//		zap.String("orderId", orderId),
	//		zap.Error(err))
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing payment success"})
	//	return
	//}

	h.logger.Info("Successfully processed payment success", zap.String("orderId", orderId))
	c.String(http.StatusOK, "Payment successful! Your booking has been confirmed.")
}

// PaymentFail handles redirect after failed payment
// GET /api/payments/fail?orderId=123
func (h *PaymentHandler) PaymentFail(c *gin.Context) {
	orderId := c.Query("orderId")
	//if orderId == "" {
	//	h.logger.Error("Missing orderId parameter in payment failure redirect")
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing orderId parameter"})
	//	return
	//}
	//
	//h.logger.Info("Payment failure redirect received", zap.String("orderId", orderId))
	//
	//err := h.paymentService.NotifyPaymentFailure(orderId)
	//if err != nil {
	//	h.logger.Error("Error processing payment failure",
	//		zap.String("orderId", orderId),
	//		zap.Error(err))
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing payment failure"})
	//	return
	//}

	h.logger.Info("Successfully processed payment failure", zap.String("orderId", orderId))
	c.String(http.StatusOK, "Payment failed. Your booking has been cancelled.")
}
