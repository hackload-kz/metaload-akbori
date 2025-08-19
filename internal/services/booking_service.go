package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BookingService interface {
	CreateBooking(req *models.CreateBookingRequest) (*models.CreateBookingResponse, error)
	GetBookingsByUser(userID int) ([]models.ListBookingsResponseItem, error)
	CancelBooking(req *models.CancelBookingRequest) error
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

func (s *bookingService) CreateBooking(req *models.CreateBookingRequest) (*models.CreateBookingResponse, error) {
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	seats, err := s.seatRepo.GetByIDs(req.SeatIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	if len(seats) != len(req.SeatIDs) {
		return nil, fmt.Errorf("some seats not found")
	}

	var totalAmount decimal.Decimal
	for _, seat := range seats {
		if seat.Status != models.SeatStatusFree {
			return nil, fmt.Errorf("seat %d is not available", seat.ID)
		}
		totalAmount = totalAmount.Add(seat.Price)
	}

	err = s.seatRepo.ReserveSeats(req.SeatIDs, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve seats: %w", err)
	}

	orderID := uuid.New().String()
	booking := &models.Booking{
		EventID:     req.EventID,
		UserID:      req.UserID,
		Status:      models.BookingStatusPending,
		TotalAmount: totalAmount,
		OrderID:     &orderID,
	}

	createdBooking, err := s.bookingRepo.Create(booking)
	if err != nil {
		s.seatRepo.ReleaseSeats(req.SeatIDs)
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return &models.CreateBookingResponse{
		ID:          createdBooking.ID,
		EventID:     createdBooking.EventID,
		UserID:      createdBooking.UserID,
		Status:      createdBooking.Status,
		TotalAmount: createdBooking.TotalAmount,
		OrderID:     createdBooking.OrderID,
		CreatedAt:   createdBooking.CreatedAt,
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

func (s *bookingService) CancelBooking(req *models.CancelBookingRequest) error {
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