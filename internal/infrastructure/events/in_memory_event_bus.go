package events

import (
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"go.uber.org/zap"
	"sync"
)

type InMemoryEventBus struct {
	subscribers map[string][]events.Handler
	mu          sync.RWMutex
	logger      *zap.SugaredLogger
}

func NewInMemoryEventBus(logger *zap.SugaredLogger) *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string][]events.Handler),
		logger:      logger,
	}
}

func (i *InMemoryEventBus) Publish(event events.Event) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for _, handler := range i.subscribers[event.GetName()] {
		go func() {
			err := handler(event, i.logger)
			if err != nil {
				i.logger.Error(err)
			}
		}()
	}

	return nil
}

func (i *InMemoryEventBus) Subscribe(event string, handler events.Handler) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.subscribers[event] = append(i.subscribers[event], handler)
	return nil
}
