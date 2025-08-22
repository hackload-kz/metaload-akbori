package models

// CreateOrderResponse представляет ответ на создание заказа
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// OrderDetails представляет детали заказа
type OrderDetails struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	StartedAt   int64  `json:"startedAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	PlacesCount int    `json:"placesCount"`
}

// Place представляет место в событии
type Place struct {
	ID     string `json:"id"`
	Row    int    `json:"row"`
	Seat   int    `json:"seat"`
	IsFree bool   `json:"is_free"`
}

// SelectPlaceRequest представляет запрос на выбор места
type SelectPlaceRequest struct {
	OrderID string `json:"orderId"`
}
