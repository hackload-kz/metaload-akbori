package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
)

type PaymentService interface {
	InitiatePayment(req *models.InitiatePaymentRequest) error
}

type paymentService struct {
	bookingRepo   repository.BookingRepository
	paymentConfig config.Payment
}

func NewPaymentService(bookingRepo repository.BookingRepository, paymentConfig config.Payment) PaymentService {
	return &paymentService{
		bookingRepo:   bookingRepo,
		paymentConfig: paymentConfig,
	}
}

func (s *paymentService) InitiatePayment(req *models.InitiatePaymentRequest) error {
	booking, err := s.bookingRepo.GetByID(req.BookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	if booking.UserID != req.UserID {
		return fmt.Errorf("unauthorized")
	}

	if booking.Status != models.BookingStatusPending {
		return fmt.Errorf("booking is not in pending status")
	}

	booking.Status = models.BookingStatusPaymentPending
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}