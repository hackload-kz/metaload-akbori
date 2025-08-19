package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"time"

	"go.uber.org/zap"
)

type PaymentGatewayService interface {
	CreatePayment(ctx context.Context, request *models.PaymentInitRequest) (*models.PaymentInitResponse, error)
	CreatePaymentRequest(orderID string, amount int64, currency, description, email string) *models.PaymentInitRequest
	CheckPaymentStatus(ctx context.Context, paymentID, orderID string) (*models.PaymentCheckResponse, error)
	ConfirmPayment(ctx context.Context, paymentID string, amount int64) (*models.PaymentConfirmResponse, error)
	CancelPayment(ctx context.Context, paymentID, reason string) (*models.PaymentCancelResponse, error)
}

type paymentGatewayService struct {
	httpClient    *http.Client
	paymentConfig config.Payment
	serviceURL    string
	logger        *zap.Logger
}

func NewPaymentGatewayService(paymentConfig config.Payment, serviceURL string, logger *zap.Logger) PaymentGatewayService {
	return &paymentGatewayService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		paymentConfig: paymentConfig,
		serviceURL:    serviceURL,
		logger:        logger,
	}
}

func (s *paymentGatewayService) CreatePayment(ctx context.Context, request *models.PaymentInitRequest) (*models.PaymentInitResponse, error) {
	s.logger.Info("Creating payment", zap.String("order_id", request.OrderID))

	// Генерируем токен аутентификации
	token := s.generateToken(request)
	request.Token = token

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/PaymentInit/init", s.paymentConfig.GatewayURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.PaymentInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Payment created successfully", zap.String("payment_id", response.PaymentID))
	return &response, nil
}

func (s *paymentGatewayService) CreatePaymentRequest(orderID string, amount int64, currency, description, email string) *models.PaymentInitRequest {
	return &models.PaymentInitRequest{
		TeamSlug:        s.paymentConfig.TeamSlug,
		OrderID:         orderID,
		Amount:          amount,
		Currency:        currency,
		Description:     description,
		Email:           email,
		Language:        "ru",
		PaymentExpiry:   3600, // 1 час
		SuccessURL:      fmt.Sprintf("%s/payments/success?orderId=%s", s.serviceURL, orderID),
		FailURL:         fmt.Sprintf("%s/payments/fail?orderId=%s", s.serviceURL, orderID),
		NotificationURL: fmt.Sprintf("%s/payments/notifications", s.serviceURL),
	}
}

func (s *paymentGatewayService) CheckPaymentStatus(ctx context.Context, paymentID, orderID string) (*models.PaymentCheckResponse, error) {
	s.logger.Info("Checking payment status", zap.String("payment_id", paymentID), zap.String("order_id", orderID))

	request := &models.PaymentCheckRequest{
		TeamSlug: s.paymentConfig.TeamSlug,
	}

	if paymentID != "" {
		request.PaymentID = paymentID
	} else if orderID != "" {
		request.OrderID = orderID
	}

	// Генерируем токен
	token := s.generateToken(request)
	request.Token = token

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/PaymentCheck/check", s.paymentConfig.GatewayURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to check payment status", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.PaymentCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Payment status checked successfully")
	return &response, nil
}

func (s *paymentGatewayService) ConfirmPayment(ctx context.Context, paymentID string, amount int64) (*models.PaymentConfirmResponse, error) {
	s.logger.Info("Confirming payment", zap.String("payment_id", paymentID), zap.Int64("amount", amount))

	request := &models.PaymentConfirmRequest{
		TeamSlug:  s.paymentConfig.TeamSlug,
		PaymentID: paymentID,
		Amount:    amount,
	}

	// Генерируем токен
	token := s.generateToken(request)
	request.Token = token

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/PaymentConfirm/confirm", s.paymentConfig.GatewayURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to confirm payment", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.PaymentConfirmResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Payment confirmed successfully", zap.String("payment_id", paymentID))
	return &response, nil
}

func (s *paymentGatewayService) CancelPayment(ctx context.Context, paymentID, reason string) (*models.PaymentCancelResponse, error) {
	s.logger.Info("Cancelling payment", zap.String("payment_id", paymentID), zap.String("reason", reason))

	request := &models.PaymentCancelRequest{
		TeamSlug:  s.paymentConfig.TeamSlug,
		PaymentID: paymentID,
		Reason:    reason,
	}

	// Генерируем токен
	token := s.generateToken(request)
	request.Token = token

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/PaymentCancel/cancel", s.paymentConfig.GatewayURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to cancel payment", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.PaymentCancelResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("Payment cancelled successfully", zap.String("payment_id", paymentID))
	return &response, nil
}

func (s *paymentGatewayService) generateToken(request interface{}) string {
	// Извлекаем поля для генерации токена в алфавитном порядке
	params := make(map[string]string)

	v := reflect.ValueOf(request).Elem()
	t := reflect.TypeOf(request).Elem()

	switch req := request.(type) {
	case *models.PaymentInitRequest:
		params["amount"] = fmt.Sprintf("%d", req.Amount)
		params["currency"] = req.Currency
		params["orderId"] = req.OrderID
		params["teamSlug"] = req.TeamSlug
	case *models.PaymentCheckRequest:
		if req.PaymentID != "" {
			params["paymentId"] = req.PaymentID
		}
		if req.OrderID != "" {
			params["orderId"] = req.OrderID
		}
		params["teamSlug"] = req.TeamSlug
	case *models.PaymentConfirmRequest:
		params["amount"] = fmt.Sprintf("%d", req.Amount)
		params["paymentId"] = req.PaymentID
		params["teamSlug"] = req.TeamSlug
	case *models.PaymentCancelRequest:
		params["paymentId"] = req.PaymentID
		params["teamSlug"] = req.TeamSlug
	default:
		// Generic reflection-based approach
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Пропускаем поле Token
			if field.Name == "Token" {
				continue
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}

			// Убираем omitempty из json тега
			fieldName := jsonTag
			if idx := len(jsonTag); idx > 0 {
				if commaIdx := len(jsonTag); commaIdx > 0 {
					for j, char := range jsonTag {
						if char == ',' {
							commaIdx = j
							break
						}
					}
					if commaIdx < len(jsonTag) {
						fieldName = jsonTag[:commaIdx]
					}
				}
			}

			if value.IsValid() && !value.IsZero() {
				params[fieldName] = fmt.Sprintf("%v", value.Interface())
			}
		}
	}

	// Сортируем ключи в алфавитном порядке
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Конкатенируем значения в алфавитном порядке ключей
	var tokenString string
	for _, key := range keys {
		tokenString += params[key]
	}

	// Добавляем пароль
	tokenString += s.paymentConfig.Password

	// Генерируем SHA-256 хеш
	hash := sha256.Sum256([]byte(tokenString))
	return fmt.Sprintf("%x", hash)
}
