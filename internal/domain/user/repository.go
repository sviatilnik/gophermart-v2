package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByLogin(ctx context.Context, login Login) (*User, error)
	Exists(ctx context.Context, login Login) (bool, error)
	Delete(ctx context.Context, id string) error
}
