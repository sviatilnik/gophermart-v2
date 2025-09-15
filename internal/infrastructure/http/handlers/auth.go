package handlers

import (
	"encoding/json"
	"github.com/sviatilnik/gophermart/internal/application/auth"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	service *auth.Service
	logger  *zap.SugaredLogger
}

func NewAuthHandler(service *auth.Service, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req auth.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Info("Login request: %v", req)

	w.Header().Add("Content-Type", "application/json")

	tokenResponse, err := h.service.LoginByPassword(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenResponse.AccessToken)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

func (h *AuthHandler) LoginByRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req auth.RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	tokenResponse, err := h.service.LoginByRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&auth.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}
