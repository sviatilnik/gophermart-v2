package middleware

import (
	"context"
	"encoding/json"
	"github.com/sviatilnik/gophermart/internal/application/auth"
	"net/http"
	"strings"
)

type userID string

var (
	RequestUserID userID
)

type AuthMiddleware struct {
	verifier auth.TokenVerifier
}

func NewAuthMiddleware(verifier auth.TokenVerifier) *AuthMiddleware {
	return &AuthMiddleware{
		verifier: verifier,
	}
}

func (m *AuthMiddleware) Handle(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.Replace(authHeader, "Bearer ", "", 1)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := m.verifier.Verify(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&auth.ErrorResponse{Error: err.Error()})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), RequestUserID, id))
		nextHandler.ServeHTTP(w, r)
	})
}
