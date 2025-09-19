package order

import (
	"context"
	"errors"
)

var (
	ErrAlreadyExists                 = errors.New("order already exists")
	ErrAlreadyCreatedByOtherCustomer = errors.New("order already created by other customer")
)

type Repository interface {
	Get(ctx context.Context, number Number) (*Order, error)
	Save(ctx context.Context, order *Order) error
	GetForCustomer(ctx context.Context, customerID string, limit uint64, offset uint64) ([]*Order, error)
	GetByStates(ctx context.Context, state []State, limit uint64, offset uint64) ([]*Order, error)
}
