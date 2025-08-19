package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"time"
)

type EventService interface {
	FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.ListEventsResponseItem, error)
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.ListEventsResponseItem, error) {
	events, err := s.eventRepo.FindEvents(query, date, page, pageSize)
	if err != nil {
		return nil, err
	}

	var response []models.ListEventsResponseItem
	for _, event := range events {
		response = append(response, models.ListEventsResponseItem{
			ID:            event.ID,
			Title:         event.Title,
			Description:   event.Description,
			Type:          event.Type,
			DatetimeStart: event.DatetimeStart,
			Provider:      event.Provider,
		})
	}

	return response, nil
}