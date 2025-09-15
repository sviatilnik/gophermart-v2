package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sviatilnik/gophermart/internal/domain/user"
)

type PostgresUserRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresUserRepository) Save(ctx context.Context, usr *user.User) error {
	exists, _ := r.Exists(ctx, usr.Login)
	if exists {
		return user.ErrUserAlreadyExists
	}

	query, _, err := r.builder.
		Insert("users").
		Columns("id", "login", "password").
		Values("?", "?", "?").
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, usr.ID, usr.Login, usr.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query, _, err := r.builder.
		Select("id", "login", "password").
		From("users").
		Where("id = ?").
		ToSql()

	if err != nil {
		return nil, err
	}

	usr := &user.User{}

	err = r.db.
		QueryRowContext(ctx, query, id).
		Scan(&usr.ID, &usr.Login, &usr.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *PostgresUserRepository) FindByLogin(ctx context.Context, login user.Login) (*user.User, error) {
	query, _, err := r.builder.
		Select("*").
		From("users").
		Where("login = ?").
		ToSql()

	if err != nil {
		return nil, err
	}

	usr := &user.User{}

	err = r.db.
		QueryRowContext(ctx, query, string(login)).
		Scan(&usr.ID, &usr.Login, &usr.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *PostgresUserRepository) Exists(ctx context.Context, login user.Login) (bool, error) {
	usr, err := r.FindByLogin(ctx, login)
	if err != nil {
		return false, err
	}

	return usr.ID != "", nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query, _, err := r.builder.Delete("users").Where("id = ?").ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
