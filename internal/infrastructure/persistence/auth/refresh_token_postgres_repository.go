package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sviatilnik/gophermart/internal/domain/auth"
	"github.com/sviatilnik/gophermart/internal/domain/user"
)

type RefreshTokenPostgresRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewRefreshTokenPostgresRepository(db *sql.DB) *RefreshTokenPostgresRepository {
	return &RefreshTokenPostgresRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *RefreshTokenPostgresRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	query, _, err := r.builder.Insert("refresh_tokens").
		Columns("user_id", "token", "expires_at").
		Values("?", "?", "?").
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenPostgresRepository) Find(ctx context.Context, token string) (*auth.RefreshToken, error) {
	query, _, err := r.builder.
		Select("token", "user_id", "expires_at").
		From("refresh_tokens").
		Where("token = ?").
		ToSql()

	if err != nil {
		return nil, err
	}

	refreshToken := &auth.RefreshToken{}

	err = r.db.
		QueryRowContext(ctx, query, token).
		Scan(&refreshToken.Token, &refreshToken.UserID, &refreshToken.ExpiresAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenPostgresRepository) Delete(ctx context.Context, token string) error {
	query, _, err := r.builder.Delete("refresh_tokens").Where("token = ?").ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return auth.ErrRefreshTokenNotFound
	}

	return nil
}
