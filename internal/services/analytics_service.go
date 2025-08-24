package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"strconv"

	"go.uber.org/zap"
)

type AnalyticsService interface {
	GetAnalytics(eventID int64) (*models.AnalyticsResponse, error)
}

type analyticsService struct {
	seatRepo    repository.SeatRepository
	bookingRepo repository.BookingRepository
	logger      *zap.Logger
}

func NewAnalyticsService(
	seatRepo repository.SeatRepository,
	bookingRepo repository.BookingRepository,
	logger *zap.Logger,
) AnalyticsService {
	return &analyticsService{
		seatRepo:    seatRepo,
		bookingRepo: bookingRepo,
		logger:      logger,
	}
}

func (s *analyticsService) GetAnalytics(eventID int64) (*models.AnalyticsResponse, error) {
	s.logger.Info("Getting analytics for event", zap.Int64("event_id", eventID))

	// Получаем статистику по местам и выручку
	seatStats, revenue, err := s.seatRepo.GetSeatStatistics(eventID)
	if err != nil {
		s.logger.Error("Failed to get seat statistics", zap.Error(err))
		return nil, err
	}

	// Получаем количество броней
	bookingsCount, _, err := s.bookingRepo.GetBookingStatistics(eventID)
	if err != nil {
		s.logger.Error("Failed to get booking statistics", zap.Error(err))
		return nil, err
	}

	// Получаем количество мест по статусам
	totalSeats := seatStats["total"]
	soldSeats := seatStats["SOLD"]
	reservedSeats := seatStats["RESERVED"]
	freeSeats := seatStats["FREE"]

	eventIDInt64, _ := strconv.ParseInt(strconv.FormatInt(eventID, 10), 10, 64)

	response := &models.AnalyticsResponse{
		EventID:       eventIDInt64,
		TotalSeats:    totalSeats,
		SoldSeats:     soldSeats,
		ReservedSeats: reservedSeats,
		FreeSeats:     freeSeats,
		TotalRevenue:  revenue,
		BookingsCount: bookingsCount,
	}

	s.logger.Info("Analytics retrieved successfully", zap.Int64("event_id", eventID))
	return response, nil
}
