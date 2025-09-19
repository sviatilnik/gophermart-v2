package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	domenevents "github.com/sviatilnik/gophermart/internal/domain/events"
)

type OutboxEvent struct {
	EventID       uuid.UUID
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       any
	Headers       map[string]any
	OccurredAt    time.Time
}

type OutboxWriter interface {
	InsertEvent(ctx context.Context, tx *sql.Tx, e OutboxEvent) error
}

type PostgresOutboxRepository struct {
	db *sql.DB
}

func NewPostgresOutboxRepository(db *sql.DB) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{db: db}
}

func (r *PostgresOutboxRepository) InsertEvent(ctx context.Context, tx *sql.Tx, e OutboxEvent) error {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}

	headersBytes, err := json.Marshal(e.Headers)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO outbox_events (
            event_id, aggregate_type, aggregate_id, event_type, payload, headers, occurred_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
    `, e.EventID, e.AggregateType, e.AggregateID, e.EventType, payloadBytes, headersBytes, e.OccurredAt)
	return err
}

// Helper to convert domain event into OutboxEvent
func ToOutboxEvent(aggregateType string, aggregateID string, event domenevents.Event, payload any, headers map[string]any) OutboxEvent {
	return OutboxEvent{
		EventID:       uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     event.GetName(),
		Payload:       payload,
		Headers:       headers,
		OccurredAt:    time.Now(),
	}
}

