package user

import "errors"

var (
	ErrLoginNotValid        = errors.New("login is not valid")
	ErrPasswordTooShort     = errors.New("password is too short (min 8 characters)")
	ErrPasswordHasNoUpper   = errors.New("password has no upper case letters")
	ErrPasswordHasNoLower   = errors.New("password has no lower case letters")
	ErrPasswordHasNoDigit   = errors.New("password has no digit letters")
	ErrPasswordHasNoSpecial = errors.New("password has no special letters")
	ErrPasswordsNotEqual    = errors.New("passwords are not equal")
)
