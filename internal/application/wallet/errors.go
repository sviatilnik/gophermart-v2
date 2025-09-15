package wallet

import "errors"

var (
	ErrNotEnoughFunds      = errors.New("not enough funds")
	ErrOrderNumberNotValid = errors.New("order number not valid")
)
