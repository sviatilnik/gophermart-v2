package auth

import (
	"context"
	"errors"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	Find(ctx context.Context, token string) (*RefreshToken, error)
	Delete(ctx context.Context, token string) error
}
