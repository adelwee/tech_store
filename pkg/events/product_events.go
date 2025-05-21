package events

type ProductCreatedEvent struct {
	ProductID uint64  `json:"product_id"`
	Name      string  `json:"name"`
	Brand     string  `json:"brand"`
	Price     float64 `json:"price"`
}

type OrderItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  uint64 `json:"quantity"`
}

type OrderCreatedEvent struct {
	OrderID uint64      `json:"order_id"`
	Items   []OrderItem `json:"items"`
}
