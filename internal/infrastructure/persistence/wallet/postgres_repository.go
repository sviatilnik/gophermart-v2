package wallet

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sviatilnik/gophermart/internal/domain/wallet"
	"time"
)

var errUnknownEvent = errors.New("unknown event")
var errVersionConflict = errors.New("version conflict")

type PostgresRepository struct {
	db                 *sql.DB
	eventsTableName    string
	withdrawsTableName string
	builder            squirrel.StatementBuilderType
}

func NewWalletPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db:                 db,
		eventsTableName:    "wallet_events",
		withdrawsTableName: "wallet_withdrawals",
		builder:            squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (p *PostgresRepository) Load(ctx context.Context, customerID string) (*wallet.Wallet, error) {
	query, _, err := p.builder.Select("event_type", "event_data").
		From(p.eventsTableName).
		Where("aggregate_id = ?").
		OrderBy("version ASC").
		ToSql()

	if err != nil {
		return nil, err
	}

	var events []wallet.Event

	rows, err := p.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var (
			eventType string
			eventData []byte
		)

		if err := rows.Scan(&eventType, &eventData); err != nil {
			return nil, err
		}

		event, err := p.unmarshalEventData(eventType, eventData)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	wlt := wallet.NewWallet(customerID)

	for _, event := range events {
		wlt.ApplyEvent(event)
	}

	wlt.Reset()

	return wlt, nil
}

func (p *PostgresRepository) Store(ctx context.Context, wlt *wallet.Wallet) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	currentVersion, err := p.checkVersion(tx, wlt)
	if err != nil {
		return err
	}

	query, _, err := p.builder.Insert(p.eventsTableName).
		Columns("event_id", "aggregate_id", "event_type", "event_data", "version", "timestamp").
		Values("?", "?", "?", "?", "?", "?").
		ToSql()
	if err != nil {
		return err
	}

	for i, event := range wlt.Events() {
		eventVersion := currentVersion + i + 1

		data, err := json.Marshal(event)
		if err != nil {
			return err
		}

		eventID := uuid.New().String()

		_, err = tx.ExecContext(
			ctx,
			query,
			eventID,
			wlt.CustomerID,
			event.GetType(),
			data,
			eventVersion,
			time.Now())
		if err != nil {
			return err
		}

		if withdrawn, ok := event.(*wallet.Withdrawn); ok {
			err = p.saveWithdrawn(tx, eventID, withdrawn)
			if err != nil {
				return err
			}
		}

	}

	// TODO add snapshots
	/* if wlt.Version() % 50 == 0 {

	} */

	return tx.Commit()
}

func (p *PostgresRepository) checkVersion(tx *sql.Tx, wlt *wallet.Wallet) (int, error) {
	var maxVersion sql.NullInt64
	query, _, err := p.builder.Select("MAX(version)").
		From(p.eventsTableName).
		Where("aggregate_id = ?").
		ToSql()
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(
		query,
		wlt.CustomerID,
	).Scan(&maxVersion)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	m := 0
	if maxVersion.Valid {
		m = int(maxVersion.Int64)
	}

	if m != wlt.Version()-len(wlt.Events()) {
		return m, errVersionConflict
	}

	return m, nil
}

func (p *PostgresRepository) saveWithdrawn(tx *sql.Tx, eventID string, withdrawn *wallet.Withdrawn) error {
	query, _, err := p.builder.Insert(p.withdrawsTableName).
		Columns("id", "event_id", "customer_id", "amount", "order_number", "timestamp").
		Values("?", "?", "?", "?", "?", "?").
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		query,
		uuid.NewString(),
		eventID,
		withdrawn.CustomerID,
		withdrawn.Amount,
		withdrawn.OrderNumber,
		withdrawn.Timestamp,
	)

	return err
}

func (p *PostgresRepository) Exists(ctx context.Context, customerID string) (bool, error) {
	query, _, err := p.builder.Select("count(*)").
		From(p.eventsTableName).
		Where("aggregate_id = ?").
		ToSql()

	if err != nil {
		return false, err
	}

	var count int
	err = p.db.QueryRowContext(ctx, query, customerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *PostgresRepository) unmarshalEventData(eventType string, data []byte) (wallet.Event, error) {
	switch eventType {
	case "created":
		var event wallet.Created
		err := json.Unmarshal(data, &event)
		return &event, err
	case "deposited":
		var event wallet.Deposited
		err := json.Unmarshal(data, &event)
		return &event, err
	case "withdrawn":
		var event wallet.Withdrawn
		err := json.Unmarshal(data, &event)
		return &event, err
	}

	return nil, errUnknownEvent
}

func (p *PostgresRepository) Withdraws(ctx context.Context, customerID string) ([]*wallet.Withdraw, error) {
	query, _, err := p.builder.Select("id", "customer_id", "amount", "order_number", "timestamp").
		From(p.withdrawsTableName).
		Where("customer_id = ?").
		OrderBy("timestamp DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	result := make([]*wallet.Withdraw, 0)
	rows, err := p.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		withdraw := &wallet.Withdraw{}

		err := rows.Scan(&withdraw.ID, &withdraw.CustomerID, &withdraw.Amount, &withdraw.OrderNumber, &withdraw.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, withdraw)
	}

	return result, nil
}
