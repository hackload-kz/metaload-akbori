package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventService interface {
	FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.ListEventsResponseItem, error)
	ClearCache()
}

type eventService struct {
	eventRepo   repository.EventRepository
	redisClient *redis.Client
	cacheTTL    time.Duration
}

func NewEventService(eventRepo repository.EventRepository, redisClient *redis.Client) EventService {
	return &eventService{
		eventRepo:   eventRepo,
		redisClient: redisClient,
		cacheTTL:    10 * time.Minute,
	}
}

func (s *eventService) generateCacheKey(query *string, date *time.Time, page, pageSize int) string {
	var queryStr string
	if query != nil {
		queryStr = *query
	}

	var dateStr string
	if date != nil {
		dateStr = date.Format("2006-01-02")
	}

	key := fmt.Sprintf("events:q:%s|d:%s|p:%d|s:%d", queryStr, dateStr, page, pageSize)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("events:%x", hash)
}

func (s *eventService) getCachedResult(ctx context.Context, cacheKey string) ([]models.ListEventsResponseItem, bool) {
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, false
	}

	var result []models.ListEventsResponseItem
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, false
	}

	return result, true
}

func (s *eventService) setCachedResult(ctx context.Context, cacheKey string, data []models.ListEventsResponseItem) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	s.redisClient.Set(ctx, cacheKey, jsonData, s.cacheTTL)
}

func (s *eventService) FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.ListEventsResponseItem, error) {
	ctx := context.Background()
	cacheKey := s.generateCacheKey(query, date, page, pageSize)

	if cachedResult, found := s.getCachedResult(ctx, cacheKey); found {
		return cachedResult, nil
	}

	events, err := s.eventRepo.FindEvents(query, date, page, pageSize)
	if err != nil {
		return nil, err
	}

	var response []models.ListEventsResponseItem
	for _, event := range events {
		response = append(response, models.ListEventsResponseItem{
			ID:    event.ID,
			Title: event.Title,
		})
	}

	s.setCachedResult(ctx, cacheKey, response)
	return response, nil
}

func (s *eventService) ClearCache() {
	ctx := context.Background()
	pattern := "events:*"

	iter := s.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}
}
