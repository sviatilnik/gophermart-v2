package auth

type TokenVerifier interface {
	Verify(token string) (userID string, err error)
}
