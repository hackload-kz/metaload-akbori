package handlers

import (
	"biletter-service/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	services *services.Services
	logger   *zap.Logger
}

func New(services *services.Services, logger *zap.Logger) *Handlers {
	return &Handlers{
		services: services,
		logger:   logger,
	}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		events := api.Group("/events")
		{
			events.GET("", h.ListEvents)
		}

		seats := api.Group("/seats")
		{
			seats.GET("/:event_id", h.ListSeats)
			seats.POST("/select", h.SelectSeat)
			seats.POST("/release", h.ReleaseSeat)
		}

		bookings := api.Group("/bookings")
		{
			bookings.POST("", h.CreateBooking)
			bookings.GET("/user/:user_id", h.ListBookingsByUser)
			bookings.POST("/cancel", h.CancelBooking)
		}

		payments := api.Group("/payments")
		{
			payments.POST("/initiate", h.InitiatePayment)
		}
	}

	router.GET("/health", h.Health)
}