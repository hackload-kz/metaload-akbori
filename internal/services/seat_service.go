package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
)

type SeatService interface {
	GetSeatsByEvent(eventID int64, page int64, pageSize int64) ([]models.ListSeatsResponseItem, error)
	SelectSeat(req *models.SelectSeatRequest) error
	ReleaseSeat(req *models.ReleaseSeatRequest) error
	FillSeats()
}

type seatService struct {
	seatRepo      repository.SeatRepository
	eventProvider EventProviderService
}

func NewSeatService(seatRepo repository.SeatRepository, eventProvider EventProviderService) SeatService {
	return &seatService{
		seatRepo:      seatRepo,
		eventProvider: eventProvider,
	}
}

func (s *seatService) GetSeatsByEvent(eventID int64, page int64, pageSize int64) ([]models.ListSeatsResponseItem, error) {
	seats, err := s.seatRepo.GetByEventID(eventID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	var response []models.ListSeatsResponseItem
	for _, seat := range seats {
		response = append(response, models.ListSeatsResponseItem{
			ID:         seat.ID,
			RowNumber:  seat.RowNumber,
			SeatNumber: seat.SeatNumber,
			Status:     seat.Status,
			Price:      seat.Price,
		})
	}

	return response, nil
}

func (s *seatService) SelectSeat(req *models.SelectSeatRequest) error {
	seat, err := s.seatRepo.GetByID(req.SeatID)
	if err != nil {
		return fmt.Errorf("failed to get seat: %w", err)
	}

	if seat == nil {
		return fmt.Errorf("seat not found")
	}

	if seat.Status != models.SeatStatusFree {
		return fmt.Errorf("seat is not available")
	}

	err = s.seatRepo.UpdateStatus(req.SeatID, models.SeatStatusReserved)
	if err != nil {
		return fmt.Errorf("failed to reserve seat: %w", err)
	}

	return nil
}

func (s *seatService) ReleaseSeat(req *models.ReleaseSeatRequest) error {
	seat, err := s.seatRepo.GetByID(req.SeatID)
	if err != nil {
		return fmt.Errorf("failed to get seat: %w", err)
	}

	if seat == nil {
		return fmt.Errorf("seat not found")
	}

	if seat.Status != models.SeatStatusReserved {
		return fmt.Errorf("seat is not reserved")
	}

	err = s.seatRepo.UpdateStatus(req.SeatID, models.SeatStatusFree)
	if err != nil {
		return fmt.Errorf("failed to release seat: %w", err)
	}

	return nil
}

func (s *seatService) FillSeats() {
	var seatCount = 1
	var eventId int64 = 1
	var page = 1
	var pageSize = 1000

	ctx := context.Background()

	for i := 0; i < 100; i++ {
		places, err := s.eventProvider.GetPlaces(ctx, &page, &pageSize)
		if err != nil {
			return
		}
		for _, place := range places {
			price := s.getPrice(seatCount)
			seat := models.Seat{
				EventID:    eventId,
				RowNumber:  place.Row,
				SeatNumber: place.Seat,
				PlaceId:    place.ID,
				Status:     models.SeatStatusFree,
				Price:      price,
			}
			err := s.seatRepo.Save(seat)
			if err != nil {
				log.Printf("failed to save seat: %v", err)
				return
			}
			seatCount++
		}
		page++
	}
}

func (s *seatService) getPrice(seatCount int) decimal.Decimal {
	if seatCount <= 10_000 {
		return decimal.NewFromInt(40_000)
	} else if seatCount <= 25_000 {
		return decimal.NewFromInt(80_000)
	} else if seatCount <= 45_000 {
		return decimal.NewFromInt(120_000)
	} else if seatCount <= 70_000 {
		return decimal.NewFromInt(160_000)
	}
	return decimal.NewFromInt(200_000)
}
