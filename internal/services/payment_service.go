package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
)

type PaymentService interface {
	InitiatePayment(req *models.InitiatePaymentRequest, userID int) (string, error)
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

func (s *paymentService) InitiatePayment(req *models.InitiatePaymentRequest, userID int) (string, error) {
	booking, err := s.bookingRepo.GetByID(req.BookingID)
	if err != nil {
		return "", fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return "", fmt.Errorf("booking not found")
	}

	// Проверяем, что пользователь является владельцем брони
	if booking.UserID != userID {
		return "", fmt.Errorf("unauthorized: booking belongs to another user")
	}

	if booking.Status != models.BookingStatusPending {
		return "", fmt.Errorf("booking is not in pending status")
	}

	booking.Status = models.BookingStatusPaymentPending
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return "", fmt.Errorf("failed to update booking: %w", err)
	}

	paymentURL := fmt.Sprintf("%s/payment?booking_id=%d", s.paymentConfig.GatewayURL, req.BookingID)
	return paymentURL, nil
}
