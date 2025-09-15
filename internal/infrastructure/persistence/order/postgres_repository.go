package order

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/sviatilnik/gophermart/internal/domain/order"
)

type PostgresRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewOrderPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) Get(ctx context.Context, number order.Number) (*order.Order, error) {
	query, _, err := r.builder.Select("id", "number", "user_id", "created_at", "state").
		From("orders").
		Where("number = ?").
		ToSql()
	if err != nil {
		return nil, err
	}

	ordr := &order.Order{}
	err = r.db.QueryRowContext(ctx, query, string(number)).
		Scan(&ordr.ID, &ordr.Number, &ordr.CustomerID, &ordr.CreatedAt, &ordr.State)
	if err != nil {
		return nil, err
	}

	return ordr, nil
}

func (r *PostgresRepository) Save(ctx context.Context, order *order.Order) error {
	query, _, err := r.builder.Insert("orders").
		Columns("id", "number", "user_id", "created_at", "state").
		Values("?", "?", "?", "?", "?").
		ToSql()

	if err != nil {
		return err
	}

	query = query + " ON CONFLICT (number) DO UPDATE SET state = $5"
	_, err = r.db.ExecContext(ctx, query, order.ID, order.Number, order.CustomerID, order.CreatedAt, order.State)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetForCustomer(ctx context.Context, customerID string, limit uint64, offset uint64) ([]*order.Order, error) {
	query, _, err := r.builder.Select("id", "number", "user_id", "created_at", "state").
		From("orders").
		Where("user_id = ?").
		OrderBy("created_at desc").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	orders := make([]*order.Order, 0)

	result, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, result.Err()
	}

	for result.Next() {
		ordr := &order.Order{}
		err := result.Scan(&ordr.ID, &ordr.Number, &ordr.CustomerID, &ordr.CreatedAt, &ordr.State)
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordr)
	}

	return orders, nil
}

func (r *PostgresRepository) GetByStates(ctx context.Context, states []order.State, limit uint64, offset uint64) ([]*order.Order, error) {
	st := make([]string, len(states))
	for i, state := range states {
		st[i] = string(state)
	}

	orders := make([]*order.Order, 0)
	result, err := r.builder.Select("id", "number", "user_id", "created_at", "state").
		From("orders").
		Where(squirrel.Eq{"state": st}).
		OrderBy("created_at desc").
		Limit(limit).
		Offset(offset).
		RunWith(r.db).
		QueryContext(ctx)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	if result.Err() != nil {
		return nil, result.Err()
	}

	for result.Next() {
		ordr := &order.Order{}
		err := result.Scan(&ordr.ID, &ordr.Number, &ordr.CustomerID, &ordr.CreatedAt, &ordr.State)
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordr)
	}

	return orders, nil
}
