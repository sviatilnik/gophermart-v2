package wallet

import "github.com/sviatilnik/gophermart/internal/domain/order"

type Command interface {
}

type DepositCommand struct {
	CustomerID string
	Amount     float64
}

func NewDepositCommand(customerID string, amount float64) *DepositCommand {
	return &DepositCommand{
		CustomerID: customerID,
		Amount:     amount,
	}
}

type WithdrawCommand struct {
	CustomerID  string
	Amount      float64
	OrderNumber order.Number
}

func NewWithdrawCommand(customerID string, orderNumber string, amount float64) (*WithdrawCommand, error) {
	n, err := order.NewOrderNumber(orderNumber)
	if err != nil {
		return nil, err
	}

	return &WithdrawCommand{
		CustomerID:  customerID,
		Amount:      amount,
		OrderNumber: n,
	}, nil
}

type CreateCommand struct {
	CustomerID string
}

func NewCreateCommand(customerID string) *CreateCommand {
	return &CreateCommand{
		CustomerID: customerID,
	}
}
