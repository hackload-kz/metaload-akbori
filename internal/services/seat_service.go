package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
)

type SeatService interface {
	GetSeatsByEvent(eventID int64) ([]models.ListSeatsResponseItem, error)
	SelectSeat(req *models.SelectSeatRequest) error
	ReleaseSeat(req *models.ReleaseSeatRequest) error
}

type seatService struct {
	seatRepo repository.SeatRepository
}

func NewSeatService(seatRepo repository.SeatRepository) SeatService {
	return &seatService{
		seatRepo: seatRepo,
	}
}

func (s *seatService) GetSeatsByEvent(eventID int64) ([]models.ListSeatsResponseItem, error) {
	seats, err := s.seatRepo.GetByEventID(eventID)
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