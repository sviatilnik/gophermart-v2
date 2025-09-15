package order

import (
	"errors"
	"unicode"
)

var (
	ErrOrderNumberNotValid = errors.New("order number not valid")
)

type Number string

func NewOrderNumber(number string) (Number, error) {
	if !isValid(number) {
		return "", ErrOrderNumberNotValid
	}

	return Number(number), nil
}

func isValid(number string) bool {
	sum := 0
	parity := len(number) % 2

	for i, char := range number {
		if !unicode.IsDigit(char) {
			return false
		}

		digit := int(char - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
