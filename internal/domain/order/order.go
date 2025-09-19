package order

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         string
	Number     Number
	CustomerID string
	CreatedAt  time.Time
	State      State
}

func NewOrder(number Number, customerID string) *Order {
	return &Order{
		ID:         uuid.NewString(),
		Number:     number,
		CustomerID: customerID,
		CreatedAt:  time.Now(),
		State:      New,
	}
}
