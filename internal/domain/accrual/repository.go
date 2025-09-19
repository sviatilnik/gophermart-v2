package accrual

import "context"

type Repository interface {
	Save(ctx context.Context, accrual *Accrual) error
	Get(ctx context.Context, orderNumber string) (*Accrual, error)
	GetForOrders(ctx context.Context, orderNumbers []string) (map[string]*Accrual, error)
}
