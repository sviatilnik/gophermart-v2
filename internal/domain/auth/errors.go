package auth

import "errors"

var (
	ErrRefreshTokenIsExpired = errors.New("refresh token is expired")
)
