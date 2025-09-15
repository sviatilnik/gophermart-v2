package wallet

import (
	"time"
)

type Event interface {
	GetType() string
}

type Created struct {
	CustomerID string
	Timestamp  time.Time
}

func (c *Created) GetType() string {
	return "created"
}

type Deposited struct {
	CustomerID string
	Amount     float64
	Timestamp  time.Time
}

func (d *Deposited) GetType() string {
	return "deposited"
}

type Withdrawn struct {
	CustomerID  string
	Amount      float64
	OrderNumber string
	Timestamp   time.Time
}

func (w *Withdrawn) GetType() string {
	return "withdrawn"
}
