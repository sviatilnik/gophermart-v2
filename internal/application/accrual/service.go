package accrual

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sviatilnik/gophermart/internal/application/order"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"go.uber.org/zap"
)

type checkOrderResponse struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Amount      float64 `json:"accrual"`
}

type Service struct {
	url          string
	repository   accrual.Repository
	orderService *order.Service
	eventBus     events.Bus
	logger       *zap.SugaredLogger
}

func NewService(url string, repository accrual.Repository, orderService *order.Service, eventBus events.Bus, logger *zap.SugaredLogger) *Service {
	return &Service{
		url:          url,
		repository:   repository,
		orderService: orderService,
		eventBus:     eventBus,
		logger:       logger,
	}
}

func (s *Service) GetAccruals(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	for {
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

			rl := NewRateLimiter()
			var wg sync.WaitGroup

			for _, o := range orders {
				wg.Add(1)
				go s.Worker(ctx, o, rl, s.url+"/api/orders/"+o.Number, &wg)
			}

			wg.Wait()
			s.logger.Info("accrual: check done")
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *Service) Worker(ctx context.Context, o *order.OrderDTO, rl *RateLimiter, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	retryCount := 5
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if rl.ShouldStop() {
				rl.WaitIfNeeded()
				continue
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if rl.HandleResponse(resp) {
				resp.Body.Close()
				rl.WaitIfNeeded()
				continue
			}

			if resp.StatusCode == http.StatusOK {
				jsonResponse := &checkOrderResponse{}
				err = json.NewDecoder(resp.Body).Decode(jsonResponse)
				if err != nil {
					return
				}

				acc := &accrual.Accrual{
					OrderNumber: jsonResponse.OrderNumber,
					State:       accrual.State(jsonResponse.Status),
					Amount:      jsonResponse.Amount,
				}

				err = s.repository.Save(ctx, acc)
				if err != nil {
					s.logger.Error("accrual: failed to save order accrual", zap.Error(err))
					continue
				}

				s.eventBus.Publish(&accrual.CreatedEvent{
					OrderNumber: acc.OrderNumber,
					Amount:      acc.Amount,
					Status:      string(acc.State),
					CustomerID:  o.CustomerID,
				})

				return
			}

			resp.Body.Close()

			retryCount--
			if retryCount > 0 {
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}
	}
}
