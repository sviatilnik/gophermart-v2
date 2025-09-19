package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID       string
	Login    Login
	Password Password
}

func NewUser(login, password string) (*User, error) {
	pass, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.NewString(),
		Login:    NewLogin(login),
		Password: pass,
	}, nil
}
