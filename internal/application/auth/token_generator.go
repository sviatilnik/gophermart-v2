package auth

import "time"

type TokenGenerator interface {
	GenerateToken(userID string, ttl time.Duration) (string, time.Time, error)
}
