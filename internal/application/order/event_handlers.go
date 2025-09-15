package order

import (
	"context"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"github.com/sviatilnik/gophermart/internal/domain/order"
	"go.uber.org/zap"
	"time"
)

func RegisterEventHandlers(bus events.Bus, orderService *Service) {
	bus.Subscribe("accrual.created", func(e events.Event, logger *zap.SugaredLogger) error {

		event := e.(*accrual.CreatedEvent)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		num, err := order.NewOrderNumber(event.OrderNumber)
		if err != nil {
			return err
		}

		o, err := orderService.orderRepo.Get(ctx, num)
		if err != nil {
			return err
		}

		o.State = order.State(event.Status)
		err = orderService.orderRepo.Save(ctx, o)
		if err != nil {
			return err
		}

		logger.Infow("new accrual created", "event", e)

		return nil
	})
}
