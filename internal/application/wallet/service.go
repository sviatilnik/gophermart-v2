package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"github.com/sviatilnik/gophermart/internal/domain/order"
	"github.com/sviatilnik/gophermart/internal/domain/wallet"
)

type Service struct {
	repo     wallet.Repository
	eventBus events.Bus
}

func NewWalletService(repo wallet.Repository, bus events.Bus) *Service {
	return &Service{
		repo:     repo,
		eventBus: bus,
	}
}

func (s *Service) Get(ctx context.Context, customerID string) (*Wallet, error) {
	wallt, err := s.repo.Load(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		CustomerID: customerID,
		Balance:    wallt.Balance,
		Withdrawn:  wallt.Withdrawn,
	}, nil
}

func (s *Service) Create(ctx context.Context, customerID string) (*Wallet, error) {
	exists, err := s.repo.Exists(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("wallet %s already exists", customerID)
	}

	w := wallet.NewWallet(customerID)

	err = w.HandleCommand(wallet.NewCreateCommand(customerID))
	if err != nil {
		return nil, err
	}

	err = s.repo.Store(ctx, w)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		CustomerID: customerID,
		Balance:    w.Balance,
		Withdrawn:  w.Withdrawn,
	}, nil
}

func (s *Service) Deposit(ctx context.Context, customerID string, amount float64) error {
	wallt, err := s.repo.Load(ctx, customerID)
	if err != nil {
		return err
	}

	err = wallt.HandleCommand(wallet.NewDepositCommand(customerID, amount))
	if err != nil {
		return err
	}

	err = s.repo.Store(ctx, wallt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, customerID string, orderNumber string, amount float64) error {
	wallt, err := s.repo.Load(ctx, customerID)
	if err != nil {
		return err
	}

	cmd, err := wallet.NewWithdrawCommand(customerID, orderNumber, amount)
	if err != nil {
		if errors.Is(err, order.ErrOrderNumberNotValid) {
			return ErrOrderNumberNotValid
		}
		return err
	}

	err = wallt.HandleCommand(cmd)
	if err != nil {
		if errors.Is(err, wallet.ErrInsufficientFunds) {
			return ErrNotEnoughFunds
		}

		return err
	}

	err = s.repo.Store(ctx, wallt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetWithdraws(ctx context.Context, customerID string) ([]*Withdraw, error) {
	withdraws, err := s.repo.Withdraws(ctx, customerID)
	if err != nil {
		return nil, err
	}

	res := make([]*Withdraw, len(withdraws))
	for i, w := range withdraws {
		res[i] = &Withdraw{
			OrderNumber: w.OrderNumber,
			Amount:      w.Amount,
			ProcessedAt: w.CreatedAt,
		}
	}

	return res, nil
}
