package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Generator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *Generator {
	return &Generator{
		secretKey: secretKey,
	}
}

func (g *Generator) GenerateToken(userID string, ttl time.Duration) (string, time.Time, error) {
	expire := time.Now().Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expire.Unix(),
	})

	signedToken, err := token.SignedString([]byte(g.secretKey))
	return signedToken, expire, err
}
