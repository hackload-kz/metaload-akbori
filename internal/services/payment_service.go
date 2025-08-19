package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"context"
	"fmt"
	"strings"
	"time"
)

type PaymentService interface {
	InitiatePayment(req *models.InitiatePaymentRequest, userID int) (string, error)
	ProcessPaymentNotification(payload *models.PaymentNotificationPayload) error
	NotifyPaymentSuccess(orderID string) error
	NotifyPaymentFailure(orderID string) error
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

func (s *paymentService) ProcessPaymentNotification(payload *models.PaymentNotificationPayload) error {
	// Поиск бронирования по paymentId или orderId
	var booking *models.Booking
	var err error

	// Сначала пытаемся найти по paymentId (если репозиторий поддерживает)
	// TODO: Добавить метод GetByPaymentID в BookingRepository

	// Пытаемся найти по orderId из data
	if payload.Data != nil {
		if orderIDRaw, exists := payload.Data["orderId"]; exists {
			orderID := fmt.Sprintf("%v", orderIDRaw)
			booking, err = s.bookingRepo.GetByOrderID(orderID)
			if err != nil {
				return fmt.Errorf("failed to get booking by order ID: %w", err)
			}
		}
	}

	if booking == nil {
		return fmt.Errorf("no booking found for payment ID: %s", payload.PaymentID)
	}

	// Обновляем paymentId если его еще нет
	if booking.PaymentID == nil {
		booking.PaymentID = &payload.PaymentID
	}

	// Обрабатываем статус платежа
	switch strings.ToUpper(payload.Status) {
	case "CONFIRMED", "COMPLETED":
		booking.Status = models.BookingStatusConfirmed
	case "FAILED", "CANCELLED", "REJECTED", "EXPIRED":
		booking.Status = models.BookingStatusCancelled
	case "AUTHORIZED":
		// Платеж авторизован, но еще не подтвержден
		// Оставляем текущий статус
	default:
		// Неизвестный статус, оставляем как есть
	}

	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}

func (s *paymentService) NotifyPaymentSuccess(orderID string) error {
	booking, err := s.bookingRepo.GetByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get booking by order ID: %w", err)
	}

	if booking == nil {
		return fmt.Errorf("booking not found for order ID: %s", orderID)
	}

	booking.Status = models.BookingStatusConfirmed
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}

func (s *paymentService) NotifyPaymentFailure(orderID string) error {
	booking, err := s.bookingRepo.GetByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get booking by order ID: %w", err)
	}

	if booking == nil {
		return fmt.Errorf("booking not found for order ID: %s", orderID)
	}

	booking.Status = models.BookingStatusCancelled
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}
