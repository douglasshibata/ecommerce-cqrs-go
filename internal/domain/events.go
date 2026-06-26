package domain

import "time"

type EventType string

const (
	ProductCreatedType  EventType = "ProductCreated"
	StockIncreasedType  EventType = "StockIncreased"
	StockDecreasedType  EventType = "StockDecreased"
	OrderPlacedType     EventType = "OrderPlaced"
)

type Event struct {
	ID            string      `json:"id"`
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Type          EventType   `json:"type"`
	Payload       interface{} `json:"payload"`
	Version       int         `json:"version"`
	CreatedAt     time.Time   `json:"created_at"`
}

type ProductCreatedPayload struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type StockIncreasedPayload struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type StockDecreasedPayload struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderPlacedPayload struct {
	OrderID    string  `json:"order_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}
