package models

// PaymentInitRequest представляет запрос на инициацию платежа
type PaymentInitRequest struct {
	TeamSlug        string `json:"teamSlug"`
	Token           string `json:"token"`
	Amount          int64  `json:"amount"`
	OrderID         string `json:"orderId"`
	Currency        string `json:"currency"`
	Description     string `json:"description"`
	SuccessURL      string `json:"successURL"`
	FailURL         string `json:"failURL"`
	NotificationURL string `json:"notificationURL"`
	PaymentExpiry   int    `json:"paymentExpiry"`
	Email           string `json:"email"`
	Language        string `json:"language"`
}

// PaymentInitResponse представляет ответ на инициацию платежа
type PaymentInitResponse struct {
	Success    *bool  `json:"success"`
	PaymentID  string `json:"paymentId"`
	OrderID    string `json:"orderId"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	PaymentURL string `json:"paymentURL"`
	ExpiresAt  string `json:"expiresAt"`
	CreatedAt  string `json:"createdAt"`
}

// PaymentCheckRequest представляет запрос на проверку статуса платежа
type PaymentCheckRequest struct {
	TeamSlug  string `json:"teamSlug"`
	Token     string `json:"token"`
	PaymentID string `json:"paymentId,omitempty"`
	OrderID   string `json:"orderId,omitempty"`
}

// PaymentCheckResponse представляет ответ на проверку статуса платежа
type PaymentCheckResponse struct {
	Success   *bool  `json:"success"`
	PaymentID string `json:"paymentId"`
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// PaymentConfirmRequest представляет запрос на подтверждение платежа
type PaymentConfirmRequest struct {
	TeamSlug  string `json:"teamSlug"`
	Token     string `json:"token"`
	PaymentID string `json:"paymentId"`
	Amount    int64  `json:"amount"`
}

// PaymentConfirmResponse представляет ответ на подтверждение платежа
type PaymentConfirmResponse struct {
	Success   *bool  `json:"success"`
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// PaymentCancelRequest представляет запрос на отмену платежа
type PaymentCancelRequest struct {
	TeamSlug  string `json:"teamSlug"`
	Token     string `json:"token"`
	PaymentID string `json:"paymentId"`
	Reason    string `json:"reason"`
}

// PaymentCancelResponse представляет ответ на отмену платежа
type PaymentCancelResponse struct {
	Success   *bool  `json:"success"`
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
