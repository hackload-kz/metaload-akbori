package handlers

import (
	"biletter-service/internal/middleware"
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
		// Публичные эндпойнты (без аутентификации)
		events := api.Group("/events")
		{
			events.GET("", h.ListEvents)
		}

		// Защищенные эндпойнты (требуют аутентификацию)
		auth := api.Group("", middleware.BasicAuth(h.services.User))
		{
			seats := auth.Group("/seats")
			{
				seats.GET("", h.ListSeats)
				seats.PATCH("/select", h.SelectSeat)
				seats.PATCH("/release", h.ReleaseSeat)
			}

			bookings := auth.Group("/bookings")
			{
				bookings.POST("", h.CreateBooking)
				bookings.GET("", h.ListBookings)
				bookings.PATCH("/initiatePayment", h.InitiatePayment)
				bookings.PATCH("/cancel", h.CancelBooking)
			}
		}
	}

	router.GET("/health", h.Health)
}
