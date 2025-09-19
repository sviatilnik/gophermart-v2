package accrual

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"time"
)

type PostgresRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (p *PostgresRepository) Save(ctx context.Context, accrual *accrual.Accrual) error {
	exists, _ := p.Get(ctx, accrual.OrderNumber)
	if exists != nil {
		return nil
	}

	query, _, err := p.builder.Insert("accruals").
		Columns("order_number", "state", "amount", "created").
		Values("?", "?", "?", "?").ToSql()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, accrual.OrderNumber, accrual.State, accrual.Amount, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) Get(ctx context.Context, orderNumber string) (*accrual.Accrual, error) {
	query, _, err := p.builder.Select("order_number", "state", "amount").
		From("accruals").
		Where("order_number = ?").
		OrderBy("created DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	a := &accrual.Accrual{}
	err = p.db.QueryRowContext(ctx, query, orderNumber).Scan(&a.OrderNumber, &a.State, &a.Amount)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (p *PostgresRepository) GetForOrders(ctx context.Context, orderNumbers []string) (map[string]*accrual.Accrual, error) {
	rows, err := p.builder.Select("order_number", "state", "amount").
		From("accruals").
		Where(squirrel.Eq{"order_number": orderNumbers}).
		OrderBy("created DESC").
		RunWith(p.db).
		QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	accruals := make(map[string]*accrual.Accrual, 0)

	for rows.Next() {
		a := &accrual.Accrual{}
		err = rows.Scan(&a.OrderNumber, &a.State, &a.Amount)
		if err != nil {
			return nil, err
		}
		accruals[a.OrderNumber] = a
	}

	return accruals, nil
}
