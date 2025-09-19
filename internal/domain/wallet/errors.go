package wallet

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrVersionConflict   = errors.New("version conflict")
)
