package wallet

import "time"

type Withdraw struct {
	ID          string
	CustomerID  string
	Amount      float64
	OrderNumber string
	CreatedAt   time.Time
}
