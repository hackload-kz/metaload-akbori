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
	CancelBooking(req *models.CancelBookingRequest, userID int) error
	SelectSeat(bookingID, seatID int64, userID int) error
	ReleaseSeat(seatID int64, userID int) error
}

type bookingService struct {
	bookingRepo     repository.BookingRepository
	bookingSeatRepo repository.BookingSeatRepository
	seatRepo        repository.SeatRepository
	eventRepo       repository.EventRepository
	txManager       *repository.TransactionManager
}

func NewBookingService(bookingRepo repository.BookingRepository, bookingSeatRepo repository.BookingSeatRepository, seatRepo repository.SeatRepository, eventRepo repository.EventRepository, txManager *repository.TransactionManager) BookingService {
	return &bookingService{
		bookingRepo:     bookingRepo,
		bookingSeatRepo: bookingSeatRepo,
		seatRepo:        seatRepo,
		eventRepo:       eventRepo,
		txManager:       txManager,
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

	var bookingIDs []int64
	for _, booking := range bookings {
		bookingIDs = append(bookingIDs, booking.ID)
	}

	// Получаем разом места по всем броням
	bookingSeats, err := s.bookingSeatRepo.GetByBookingIDs(bookingIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking seats: %w", err)
	}

	var response []models.ListBookingsResponseItem
	for _, booking := range bookings {
		// Заполняем места
		seats := []models.ListBookingsResponseItemSeat{}
		for _, bookingSeat := range bookingSeats {
			if bookingSeat.BookingID == booking.ID {
				seats = append(seats, models.ListBookingsResponseItemSeat{SeatID: bookingSeat.SeatID})
			}
		}

		item := models.ListBookingsResponseItem{
			ID:      booking.ID,
			EventID: booking.EventID,
			Seats:   seats,
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
	return s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
		// Используем пессимистичную блокировку для места
		seat, err := txRepo.Seat.GetByIDForUpdate(seatID)
		if err != nil {
			return fmt.Errorf("failed to get seat for update: %w", err)
		}
		if seat == nil {
			return fmt.Errorf("seat not found")
		}

		// Проверяем, что место свободно
		if seat.Status != models.SeatStatusFree {
			return fmt.Errorf("seat is not available")
		}

		// Получаем бронь
		booking, err := txRepo.Booking.GetByID(bookingID)
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

		// Резервируем место
		seat.Status = models.SeatStatusReserved
		err = txRepo.Seat.Update(seat)
		if err != nil {
			return fmt.Errorf("failed to reserve seat: %w", err)
		}

		// Создаем связь брони с местом
		bookingSeat := &models.BookingSeat{BookingID: bookingID, SeatID: seatID}
		_, err = txRepo.BookingSeat.Create(bookingSeat)
		if err != nil {
			return fmt.Errorf("failed to create booking seat: %w", err)
		}

		return nil
	})
}

func (s *bookingService) ReleaseSeat(seatID int64, userID int) error {
	return s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
		seat, err := txRepo.Seat.GetByIDForUpdate(seatID)
		if err != nil {
			return fmt.Errorf("failed to get seat: %w", err)
		}
		if seat == nil {
			return fmt.Errorf("seat not found")
		}
		if seat.Status != models.SeatStatusReserved {
			return fmt.Errorf("seat is not reserved")
		}

		err = txRepo.Seat.UpdateStatus(seatID, models.SeatStatusFree)
		if err != nil {
			return fmt.Errorf("failed to release seat: %w", err)
		}

		// Удаляем связь места с бронью
		bookingSeats, err := txRepo.BookingSeat.GetBySeatID(seatID)
		if err != nil {
			return fmt.Errorf("failed to get booking seats: %w", err)
		}
		for _, bookingSeat := range bookingSeats {
			err = txRepo.BookingSeat.Delete(bookingSeat.ID)
			if err != nil {
				return fmt.Errorf("failed to delete booking seat: %w", err)
			}
		}

		return nil
	})
}
