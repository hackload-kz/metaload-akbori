package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BookingService interface {
	CreateBooking(req *models.CreateBookingRequest, userID int) (*models.CreateBookingResponse, error)
	GetBookingsByUser(userID int) ([]models.ListBookingsResponseItem, error)
	GetAllBookings() ([]models.ListBookingsResponseItem, error)
	CancelBooking(req *models.CancelBookingRequest, userID int) error
	SelectSeat(bookingID, seatID int64, userID int) error
	ReleaseSeat(seatID int64, userID int) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	seatRepo    repository.SeatRepository
	eventRepo   repository.EventRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, seatRepo repository.SeatRepository, eventRepo repository.EventRepository) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		seatRepo:    seatRepo,
		eventRepo:   eventRepo,
	}
}

func (s *bookingService) CreateBooking(req *models.CreateBookingRequest, userID int) (*models.CreateBookingResponse, error) {
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	orderID := uuid.New().String()
	booking := &models.Booking{
		EventID:     req.EventID,
		UserID:      userID,
		Status:      models.BookingStatusPending,
		TotalAmount: decimal.Zero,
		OrderID:     &orderID,
	}

	createdBooking, err := s.bookingRepo.Create(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return &models.CreateBookingResponse{
		ID: createdBooking.ID,
	}, nil
}

func (s *bookingService) GetBookingsByUser(userID int) ([]models.ListBookingsResponseItem, error) {
	bookings, err := s.bookingRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	var response []models.ListBookingsResponseItem
	for _, booking := range bookings {
		event, err := s.eventRepo.GetByID(booking.EventID)
		if err != nil {
			return nil, fmt.Errorf("failed to get event: %w", err)
		}

		item := models.ListBookingsResponseItem{
			ID:          booking.ID,
			EventTitle:  event.Title,
			Status:      booking.Status,
			TotalAmount: booking.TotalAmount,
			PaymentID:   booking.PaymentID,
			OrderID:     booking.OrderID,
			CreatedAt:   booking.CreatedAt,
			Seats:       []models.ListBookingsResponseItemSeat{},
		}

		response = append(response, item)
	}

	return response, nil
}

func (s *bookingService) GetAllBookings() ([]models.ListBookingsResponseItem, error) {
	bookings, err := s.bookingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	var response []models.ListBookingsResponseItem
	for _, booking := range bookings {
		event, err := s.eventRepo.GetByID(booking.EventID)
		if err != nil {
			return nil, fmt.Errorf("failed to get event: %w", err)
		}

		item := models.ListBookingsResponseItem{
			ID:          booking.ID,
			EventTitle:  event.Title,
			Status:      booking.Status,
			TotalAmount: booking.TotalAmount,
			PaymentID:   booking.PaymentID,
			OrderID:     booking.OrderID,
			CreatedAt:   booking.CreatedAt,
			Seats:       []models.ListBookingsResponseItemSeat{},
		}

		response = append(response, item)
	}

	return response, nil
}

func (s *bookingService) CancelBooking(req *models.CancelBookingRequest, userID int) error {
	booking, err := s.bookingRepo.GetByID(req.BookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	// Проверяем, что пользователь является владельцем брони
	if booking.UserID != userID {
		return fmt.Errorf("unauthorized: booking belongs to another user")
	}

	if booking.Status == models.BookingStatusCancelled {
		return fmt.Errorf("booking already cancelled")
	}

	booking.Status = models.BookingStatusCancelled
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}

func (s *bookingService) SelectSeat(bookingID, seatID int64, userID int) error {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	// Проверяем, что пользователь является владельцем брони
	if booking.UserID != userID {
		return fmt.Errorf("unauthorized: booking belongs to another user")
	}

	seat, err := s.seatRepo.GetByID(seatID)
	if err != nil {
		return fmt.Errorf("failed to get seat: %w", err)
	}
	if seat == nil {
		return fmt.Errorf("seat not found")
	}

	if seat.Status != models.SeatStatusFree {
		return fmt.Errorf("seat is not available")
	}

	err = s.seatRepo.UpdateStatus(seatID, models.SeatStatusReserved)
	if err != nil {
		return fmt.Errorf("failed to reserve seat: %w", err)
	}

	return nil
}

func (s *bookingService) ReleaseSeat(seatID int64, userID int) error {
	seat, err := s.seatRepo.GetByID(seatID)
	if err != nil {
		return fmt.Errorf("failed to get seat: %w", err)
	}
	if seat == nil {
		return fmt.Errorf("seat not found")
	}

	if seat.Status != models.SeatStatusReserved {
		return fmt.Errorf("seat is not reserved")
	}

	// Находим бронь, к которой привязано место
	// В Go версии нет BookingSeat таблицы, поэтому пока просто освобождаем место
	// TODO: Добавить проверку владельца места через BookingSeat таблицу

	err = s.seatRepo.UpdateStatus(seatID, models.SeatStatusFree)
	if err != nil {
		return fmt.Errorf("failed to release seat: %w", err)
	}

	return nil
}
