package accrual

import (
	"context"
	"github.com/sviatilnik/gophermart/internal/application/order"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	service      accrual.Service
	repository   accrual.Repository
	orderService *order.Service
	eventBus     events.Bus
	logger       *zap.SugaredLogger
}

func NewService(service accrual.Service, repository accrual.Repository, orderService *order.Service, eventBus events.Bus, logger *zap.SugaredLogger) *Service {
	return &Service{
		service:      service,
		repository:   repository,
		orderService: orderService,
		eventBus:     eventBus,
		logger:       logger,
	}
}

func (s *Service) GetAccruals(ctx context.Context) {
	for {
		ticker := time.NewTicker(15 * time.Second)

		select {
		case <-ctx.Done():
			s.logger.Info("accrual: service shutting down")
		case <-ticker.C:
			s.logger.Info("accrual: starting ...")

			orders, err := s.orderService.GetUnprocessedOrders(ctx, 100, 0)
			if err != nil {
				s.logger.Error("accrual: failed to fetch unprocessed orders", zap.Error(err))
				continue
			}

			if len(orders) == 0 {
				s.logger.Info("accrual: no unprocessed orders")
				continue
			}

			for _, o := range orders {
				checkOrder, err := s.service.CheckOrder(o.Number)
				if err != nil {
					s.logger.Error("accrual: failed to check order", zap.Error(err))
					continue
				}

				err = s.repository.Save(ctx, checkOrder)
				if err != nil {
					s.logger.Error("accrual: failed to save order accrual", zap.Error(err))
					continue
				}

				s.eventBus.Publish(&accrual.CreatedEvent{
					OrderNumber: checkOrder.OrderNumber,
					Amount:      checkOrder.Amount,
					Status:      string(checkOrder.State),
					CustomerID:  o.CustomerID,
				})
			}
		}
	}
}
