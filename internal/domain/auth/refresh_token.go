package auth

import "time"

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

func (t *RefreshToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t *RefreshToken) CheckExpiration() error {
	if t.IsExpired() {
		return ErrRefreshTokenIsExpired
	}

	return nil
}
