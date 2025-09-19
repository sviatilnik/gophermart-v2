package user

import "golang.org/x/crypto/bcrypt"

type Password string

func NewPassword(password string) (Password, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return Password(hashedPassword), nil
}

func (p Password) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	if err != nil {
		return ErrPasswordsNotEqual
	}

	return nil
}
