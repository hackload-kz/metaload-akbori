package services

import (
	"biletter-service/internal/repository"
	"fmt"

	"go.uber.org/zap"
)

type ResetService interface {
	ResetAllData() error
}

type resetService struct {
	bookingRepo repository.BookingRepository
	seatRepo    repository.SeatRepository
	txManager   *repository.TransactionManager
	logger      *zap.Logger
}

func NewResetService(
	bookingRepo repository.BookingRepository,
	seatRepo repository.SeatRepository,
	txManager *repository.TransactionManager,
	logger *zap.Logger,
) ResetService {
	return &resetService{
		bookingRepo: bookingRepo,
		seatRepo:    seatRepo,
		txManager:   txManager,
		logger:      logger,
	}
}

func (s *resetService) ResetAllData() error {
	s.logger.Info("Starting data reset")

	// Выполняем все операции в одной транзакции
	return s.txManager.WithTransaction(func(repos *repository.TransactionRepository) error {
		// 1. Удаляем все брони и связанные места
		if err := repos.Booking.DeleteAll(); err != nil {
			s.logger.Error("Failed to delete bookings", zap.Error(err))
			return fmt.Errorf("failed to delete bookings: %w", err)
		}
		s.logger.Info("All bookings deleted")

		// 2. Сбрасываем статус всех мест на FREE
		if err := repos.Seat.ResetAllStatus(); err != nil {
			s.logger.Error("Failed to reset seats status", zap.Error(err))
			return fmt.Errorf("failed to reset seats status: %w", err)
		}
		s.logger.Info("All seats status reset to FREE")

		s.logger.Info("Data reset completed successfully")
		return nil
	})
}
