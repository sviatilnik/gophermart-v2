package auth

import "time"

type AccessToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}
