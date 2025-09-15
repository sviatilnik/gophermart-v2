package wallet

import "time"

type Wallet struct {
	CustomerID string  `json:"-"`
	Balance    float64 `json:"current"`
	Withdrawn  float64 `json:"withdrawn"`
}

type CreateWithdrawal struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
}

type Withdraw struct {
	OrderNumber string    `json:"order"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
