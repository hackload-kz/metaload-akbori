package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type EventProviderService interface {
	CreateOrder(ctx context.Context) (*models.CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderID string) (*models.OrderDetails, error)
	SubmitOrder(ctx context.Context, orderID string) error
	ConfirmOrder(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
	GetPlaces(ctx context.Context, page, pageSize *int) ([]*models.Place, error)
	GetPlace(ctx context.Context, placeID string) (*models.Place, error)
	SelectPlace(ctx context.Context, placeID, orderID string) error
	ReleasePlace(ctx context.Context, placeID string) error
}

type eventProviderService struct {
	httpClient *http.Client
	config     config.ExternalService
	logger     *zap.Logger
}

func NewEventProviderService(cfg config.ExternalService, logger *zap.Logger) EventProviderService {
	return &eventProviderService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
		logger: logger,
	}
}

func (s *eventProviderService) CreateOrder(ctx context.Context) (*models.CreateOrderResponse, error) {
	s.logger.Info("Creating order in Hackload service")

	startTime := time.Now()
	defer func() {
		s.logger.Info("Order created", zap.Duration("exec_time", time.Since(startTime)))
	}()

	url := fmt.Sprintf("%s/api/partners/%s/orders", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Order created successfully", zap.String("order_id", response.OrderID))
	return &response, nil
}

func (s *eventProviderService) GetOrder(ctx context.Context, orderID string) (*models.OrderDetails, error) {
	s.logger.Info("Getting order details", zap.String("order_id", orderID))

	url := fmt.Sprintf("%s/api/partners/%s/orders/%s", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, orderID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to get order details", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.OrderDetails
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Order details retrieved", zap.String("status", response.Status))
	return &response, nil
}

func (s *eventProviderService) SubmitOrder(ctx context.Context, orderID string) error {
	s.logger.Info("Submitting order", zap.String("order_id", orderID))

	url := fmt.Sprintf("%s/api/partners/%s/orders/%s/submit", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, orderID)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to submit order", zap.Error(err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Order submitted successfully", zap.String("order_id", orderID))
	return nil
}

func (s *eventProviderService) ConfirmOrder(ctx context.Context, orderID string) error {
	s.logger.Info("Confirming order", zap.String("order_id", orderID))

	url := fmt.Sprintf("%s/api/partners/%s/orders/%s/confirm", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, orderID)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to confirm order", zap.Error(err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Order confirmed successfully", zap.String("order_id", orderID))
	return nil
}

func (s *eventProviderService) CancelOrder(ctx context.Context, orderID string) error {
	s.logger.Info("Cancelling order", zap.String("order_id", orderID))

	url := fmt.Sprintf("%s/api/partners/%s/orders/%s/cancel", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, orderID)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to cancel order", zap.Error(err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Order cancelled successfully", zap.String("order_id", orderID))
	return nil
}

func (s *eventProviderService) GetPlaces(ctx context.Context, page, pageSize *int) ([]*models.Place, error) {
	pageVal := 1
	pageSizeVal := 20
	if page != nil {
		pageVal = *page
	}
	if pageSize != nil {
		pageSizeVal = *pageSize
	}

	s.logger.Info("Getting places", zap.Int("page", pageVal), zap.Int("page_size", pageSizeVal))

	url := fmt.Sprintf("%s/api/partners/%s/places?page=%d&pageSize=%d",
		s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, pageVal, pageSizeVal)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to get places", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var places []*models.Place
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Places retrieved", zap.Int("count", len(places)))
	return places, nil
}

func (s *eventProviderService) GetPlace(ctx context.Context, placeID string) (*models.Place, error) {
	s.logger.Info("Getting place details", zap.String("place_id", placeID))

	url := fmt.Sprintf("%s/api/partners/%s/places/%s", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, placeID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to get place details", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var place models.Place
	if err := json.NewDecoder(resp.Body).Decode(&place); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Place details retrieved", zap.String("place_id", place.ID))
	return &place, nil
}

func (s *eventProviderService) SelectPlace(ctx context.Context, placeID, orderID string) error {
	s.logger.Info("Selecting place", zap.String("place_id", placeID), zap.String("order_id", orderID))

	request := models.SelectPlaceRequest{
		OrderID: orderID,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/partners/%s/places/%s/select", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, placeID)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to select place", zap.Error(err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Place selected successfully", zap.String("place_id", placeID), zap.String("order_id", orderID))
	return nil
}

func (s *eventProviderService) ReleasePlace(ctx context.Context, placeID string) error {
	s.logger.Info("Releasing place", zap.String("place_id", placeID))

	url := fmt.Sprintf("%s/api/partners/%s/places/%s/release", s.config.Hackload.BaseURL, s.config.Hackload.APIVersion, placeID)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to release place", zap.Error(err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Place released successfully", zap.String("place_id", placeID))
	return nil
}
