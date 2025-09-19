package user

import "unicode"

type PasswordChecker interface {
	Check(password string) error
}

type SimplePasswordChecker struct{}

func NewSimplePasswordChecker() *SimplePasswordChecker {
	return &SimplePasswordChecker{}
}

func (s *SimplePasswordChecker) Check(password string) error {
	if len(password) == 0 {
		return ErrPasswordTooShort
	}

	return nil
}

type PasswordCheckerService struct{}

func NewPasswordCheckerService() *PasswordCheckerService {
	return &PasswordCheckerService{}
}

func (p *PasswordCheckerService) Check(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordHasNoUpper
	}
	if !hasLower {
		return ErrPasswordHasNoLower
	}
	if !hasDigit {
		return ErrPasswordHasNoDigit
	}
	if !hasSpecial {
		return ErrPasswordHasNoSpecial
	}

	return nil
}
