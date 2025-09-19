package events

import (
	"context"
	"database/sql"
	"time"

	domenevents "github.com/sviatilnik/gophermart/internal/domain/events"
	"go.uber.org/zap"
)

// OutboxBus реализует шину событий, которая при Publish пишет событие в outbox.
// Subscribe делегирует во внутреннюю in-memory шину, чтобы обработчики были зарегистрированы.
type OutboxBus struct {
	db     *sql.DB
	inner  domenevents.Bus
	writer OutboxWriter
	logger *zap.SugaredLogger
}

func NewOutboxBus(db *sql.DB, inner domenevents.Bus, logger *zap.SugaredLogger, writer OutboxWriter) *OutboxBus {
	return &OutboxBus{db: db, inner: inner, writer: writer, logger: logger}
}

func (b *OutboxBus) Publish(event domenevents.Event) error {
	// Пишем в outbox вне внешней транзакции (упрощённый вариант).
	// Диспетчер доставит до inner позже.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	out := OutboxEvent{
		AggregateType: "system",
		AggregateID:   "",
		EventType:     event.GetName(),
		Payload:       event,
		Headers:       map[string]any{},
		OccurredAt:    time.Now(),
	}

	if err := b.writer.InsertEvent(ctx, tx, out); err != nil {
		return err
	}

	return tx.Commit()
}

func (b *OutboxBus) Subscribe(event string, handler domenevents.Handler) error {
	return b.inner.Subscribe(event, handler)
}
