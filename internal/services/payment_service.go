package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"context"
	"fmt"
	"time"
)

type PaymentService interface {
	InitiatePayment(req *models.InitiatePaymentRequest, userID int) (string, error)
}

type paymentService struct {
	bookingRepo           repository.BookingRepository
	paymentConfig         config.Payment
	paymentGatewayService PaymentGatewayService
	userService           UserService
}

func NewPaymentService(bookingRepo repository.BookingRepository, paymentConfig config.Payment, paymentGatewayService PaymentGatewayService, userService UserService) PaymentService {
	return &paymentService{
		bookingRepo:           bookingRepo,
		paymentConfig:         paymentConfig,
		paymentGatewayService: paymentGatewayService,
		userService:           userService,
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

	// Получаем пользователя для email
	user, err := s.userService.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Обновляем статус на PAYMENT_PENDING
	booking.Status = models.BookingStatusPaymentPending
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return "", fmt.Errorf("failed to update booking: %w", err)
	}

	// Сумма в тыйынах (умножаем на 100)
	amountInTiyn := booking.TotalAmount.IntPart() * 100

	// Создаем запрос на платеж
	paymentRequest := s.paymentGatewayService.CreatePaymentRequest(
		*booking.OrderID,
		amountInTiyn,
		"KZT",
		fmt.Sprintf("Оплата бронирования #%s", *booking.OrderID),
		user.Email,
	)

	// Создаем платеж в платежном шлюзе
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	paymentResponse, err := s.paymentGatewayService.CreatePayment(ctx, paymentRequest)
	if err != nil {
		// В случае ошибки возвращаем бронирование в исходное состояние
		booking.Status = models.BookingStatusPending
		s.bookingRepo.Update(booking)
		return "", fmt.Errorf("failed to create payment: %w", err)
	}

	return paymentResponse.PaymentURL, nil
}
