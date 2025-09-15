package wallet

import (
	"context"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"github.com/sviatilnik/gophermart/internal/domain/user"
	"go.uber.org/zap"
	"time"
)

func RegisterEventHandlers(bus events.Bus, walletService *Service) {
	bus.Subscribe("user.registered", func(e events.Event, logger *zap.SugaredLogger) error {

		event := e.(*user.Registered)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := walletService.Create(ctx, event.UserID)
		if err != nil {
			logger.Error("wallet creation failed", zap.Error(err))
			return err
		}

		return nil
	})

	bus.Subscribe("accrual.created", func(e events.Event, logger *zap.SugaredLogger) error {
		event := e.(*accrual.CreatedEvent)

		if event.Status == string(accrual.Processed) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := walletService.Deposit(ctx, event.CustomerID, event.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
