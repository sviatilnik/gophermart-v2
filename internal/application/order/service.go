package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"github.com/sviatilnik/gophermart/internal/domain/order"
	"github.com/sviatilnik/gophermart/internal/domain/user"
	"time"
)

type Service struct {
	orderRepo   order.Repository
	userRepo    user.Repository
	accrualRepo accrual.Repository
}

func NewOrderService(orderRepo order.Repository, userRepo user.Repository, accrualRepo accrual.Repository) *Service {
	return &Service{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		accrualRepo: accrualRepo,
	}
}

func (s *Service) Create(ctx context.Context, req CreateOrderDTO) (*OrderDTO, error) {
	orderNumber, err := order.NewOrderNumber(req.Number)
	if err != nil {
		return nil, err
	}

	usr, err := s.userRepo.FindByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}

	existsOrder, err := s.orderRepo.Get(ctx, orderNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existsOrder != nil {
		if existsOrder.CustomerID != req.CustomerID {
			return nil, ErrAlreadyCreatedByOtherCustomer
		}

		return nil, ErrAlreadyExists
	}

	newOrder := order.NewOrder(orderNumber, usr.ID)

	err = s.orderRepo.Save(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	return &OrderDTO{
		OrderID:    newOrder.ID,
		Status:     string(newOrder.State),
		Number:     string(newOrder.Number),
		UploadedAt: newOrder.CreatedAt.UTC().Format(time.RFC3339),
		CustomerID: newOrder.CustomerID,
	}, nil
}

func (s *Service) GetOrders(ctx context.Context, customerID string, limit int, offset int) ([]*OrderDTO, error) {
	customerOrders, err := s.orderRepo.GetForCustomer(ctx, customerID, uint64(limit), uint64(offset))
	if err != nil {
		return nil, err
	}

	oNumbers := make([]string, len(customerOrders))
	for i, o := range customerOrders {
		oNumbers[i] = string(o.Number)
	}

	acc, err := s.accrualRepo.GetForOrders(ctx, oNumbers)
	if err != nil {
		return nil, err
	}

	result := make([]*OrderDTO, 0)
	for _, ordr := range customerOrders {
		n := string(ordr.Number)

		o := &OrderDTO{
			OrderID:    ordr.ID,
			Status:     string(ordr.State),
			Number:     n,
			UploadedAt: ordr.CreatedAt.UTC().Format(time.RFC3339),
			CustomerID: ordr.CustomerID,
		}

		a, has := acc[n]
		if has {
			o.Accrual = a.Amount
		}
		result = append(result, o)
	}

	return result, nil
}

func (s *Service) GetUnprocessedOrders(ctx context.Context, limit int, offset int) ([]*OrderDTO, error) {
	customerOrders, err := s.orderRepo.GetByStates(ctx, []order.State{order.New, order.Processing}, uint64(limit), uint64(offset))
	if err != nil {
		return nil, err
	}

	result := make([]*OrderDTO, 0)
	for _, o := range customerOrders {
		result = append(result, &OrderDTO{
			OrderID:    o.ID,
			Status:     string(o.State),
			Number:     string(o.Number),
			UploadedAt: o.CreatedAt.UTC().Format(time.RFC3339),
			CustomerID: o.CustomerID,
		})
	}

	return result, nil
}
