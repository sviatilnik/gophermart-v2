package wallet

import "context"

type Repository interface {
	Load(ctx context.Context, customerID string) (*Wallet, error)
	Store(ctx context.Context, wallet *Wallet) error
	Exists(ctx context.Context, customerID string) (bool, error)
	Withdraws(ctx context.Context, customerID string) ([]*Withdraw, error)
}
