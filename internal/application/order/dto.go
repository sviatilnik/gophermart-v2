package order

// CreateOrderDTO - DTO для входящих данных
type CreateOrderDTO struct {
	Number     string `json:"number"`
	CustomerID string `json:"customer_id"`
}

// OrderDTO - DTO для ответа
type OrderDTO struct {
	OrderID    string  `json:"order_id"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	Number     string  `json:"number"`
	UploadedAt string  `json:"uploaded_at"`
	CustomerID string  `json:"-"`
}
