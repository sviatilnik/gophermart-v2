package user

import (
	"strings"
)

type Login string

func NewLogin(login string) Login {
	return Login(strings.TrimSpace(login))
}
