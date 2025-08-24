package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"biletter-service/pkg/broker"
	"context"
	"fmt"
	"strconv"

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
	eventPublisher  broker.Publisher
	bookingTopic    string
}

func NewBookingService(bookingRepo repository.BookingRepository, bookingSeatRepo repository.BookingSeatRepository, seatRepo repository.SeatRepository, eventRepo repository.EventRepository, txManager *repository.TransactionManager, eventPublisher broker.Publisher, bookingTopic string) BookingService {
	return &bookingService{
		bookingRepo:     bookingRepo,
		bookingSeatRepo: bookingSeatRepo,
		seatRepo:        seatRepo,
		eventRepo:       eventRepo,
		txManager:       txManager,
		eventPublisher:  eventPublisher,
		bookingTopic:    bookingTopic,
	}
}

func (s *bookingService) CreateBooking(req *models.CreateBookingRequest, userID int) (*models.CreateBookingResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Сохраняем для отправки событии
	var createdBooking *models.Booking

	err := s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
		event, err := s.eventRepo.GetByID(req.EventID)
		if err != nil {
			return fmt.Errorf("failed to get event: %w", err)
		}
		if event == nil {
			return fmt.Errorf("event not found")
		}

		orderID := uuid.New().String()
		booking := &models.Booking{
			EventID:     req.EventID,
			UserID:      userID,
			Status:      models.BookingStatusPending,
			TotalAmount: decimal.Zero,
			OrderID:     &orderID,
		}

		created, err := s.bookingRepo.Create(booking)
		if err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}
		createdBooking = created

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Отправляем событие создания брони
	eventData := models.BookingCreatedData{
		BookingID:   createdBooking.ID,
		EventID:     createdBooking.EventID,
		UserID:      createdBooking.UserID,
		TotalAmount: createdBooking.TotalAmount.Mul(decimal.NewFromInt(100)).IntPart(),
	}
	s.publishEvent(models.BookingCreatedEvent, createdBooking.ID, eventData)

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

	var response = []models.ListBookingsResponseItem{}
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
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// Сохраняем для отправки событии
	var removedBookingSeats []*models.BookingSeat

	err := s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
		booking, err := txRepo.Booking.GetByIDForUpdate(req.BookingID)
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
		err = txRepo.Booking.Update(booking)
		if err != nil {
			return fmt.Errorf("failed to update booking: %w", err)
		}

		// Получаем места
		bookingSeats, err := txRepo.BookingSeat.GetByBookingID(booking.ID)
		if err != nil {
			return fmt.Errorf("failed to get booking seats: %w", err)
		}
		removedBookingSeats = bookingSeats

		for _, bookingSeat := range bookingSeats {
			err = txRepo.BookingSeat.Delete(bookingSeat.ID)
			if err != nil {
				return fmt.Errorf("failed to delete booking seat: %w", err)
			}
		}

		return nil
	})

	if err == nil {
		if removedBookingSeats != nil && len(removedBookingSeats) > 0 {
			for _, bookingSeat := range removedBookingSeats {
				// Отправляем событие освобождения места
				eventData := models.SeatReleasedData{
					BookingID: req.BookingID,
					SeatID:    bookingSeat.SeatID,
					UserID:    userID,
				}
				s.publishEvent(models.SeatReleasedEvent, req.BookingID, eventData)
			}
		}

		// Отправляем событие отмены брони
		eventData := models.BookingCancelledData{
			BookingID: req.BookingID,
			UserID:    userID,
			Reason:    "cancelled_by_user",
		}
		s.publishEvent(models.BookingCancelledEvent, req.BookingID, eventData)
	}

	return err
}

func (s *bookingService) SelectSeat(bookingID, seatID int64, userID int) error {
	err := s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
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

	// Отправляем событие выбора места после успешного завершения транзакции
	if err == nil {
		eventData := models.SeatSelectedData{
			BookingID: bookingID,
			SeatID:    seatID,
			UserID:    userID,
		}
		s.publishEvent(models.SeatSelectedEvent, bookingID, eventData)
	}

	return err
}

func (s *bookingService) ReleaseSeat(seatID int64, userID int) error {
	// Сохраняем для отправки события
	var releasedBookingID int64

	err := s.txManager.WithTransaction(func(txRepo *repository.TransactionRepository) error {
		// Получаем мести с пессимистичной блокировкой
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

		// Получаем связи места с бронью до удаления
		bookingSeats, err := txRepo.BookingSeat.GetBySeatID(seatID)
		if err != nil {
			return fmt.Errorf("failed to get booking seats: %w", err)
		}

		if len(bookingSeats) > 0 {
			releasedBookingID = bookingSeats[0].BookingID
			booking, err := txRepo.Booking.GetByID(releasedBookingID)
			if err != nil {
				return fmt.Errorf("failed to get booking: %w", err)
			}
			if booking.UserID != userID {
				return fmt.Errorf("unauthorized: booking belongs to another user")
			}
		}

		err = txRepo.Seat.UpdateStatus(seatID, models.SeatStatusFree)
		if err != nil {
			return fmt.Errorf("failed to release seat: %w", err)
		}

		// Удаляем связь места с бронью и сохраняем ID брони для события
		for _, bookingSeat := range bookingSeats {
			err = txRepo.BookingSeat.Delete(bookingSeat.ID)
			if err != nil {
				return fmt.Errorf("failed to delete booking seat: %w", err)
			}
		}

		return nil
	})

	// Отправляем событие освобождения места после успешного завершения транзакции
	if err == nil && releasedBookingID > 0 {
		eventData := models.SeatReleasedData{
			BookingID: releasedBookingID,
			SeatID:    seatID,
			UserID:    userID,
		}
		s.publishEvent(models.SeatReleasedEvent, releasedBookingID, eventData)
	}

	return err
}

// publishEvent отправляет событие в Broker с ID брони в качестве ключа
func (s *bookingService) publishEvent(eventType models.EventType, bookingID int64, data any) {
	if s.eventPublisher == nil {
		return // Graceful degradation если publisher не настроен
	}

	bookingIDStr := strconv.FormatInt(bookingID, 10)
	event := models.NewDomainEvent(eventType, bookingIDStr, data)

	ctx := context.Background()
	if err := s.eventPublisher.Publish(ctx, s.bookingTopic, bookingIDStr, event); err != nil {
		// В продакшене здесь должно быть логирование ошибки
		fmt.Printf("Failed to publish event %s into topic %s for booking %d: %v\n", eventType, s.bookingTopic, bookingID, err)
	}
}
