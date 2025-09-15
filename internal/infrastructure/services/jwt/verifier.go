package jwt

import "github.com/golang-jwt/jwt/v4"

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
	}
}

func (v *Verifier) Verify(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	tok, err := jwt.ParseWithClaims(
		token,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(v.secret), nil
		})
	if err != nil {
		return "", err
	}

	if !tok.Valid {
		return "", nil
	}

	return claims.Subject, nil
}
