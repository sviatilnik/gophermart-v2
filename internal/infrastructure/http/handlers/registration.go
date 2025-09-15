package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/gophermart/internal/application/auth"
	"github.com/sviatilnik/gophermart/internal/domain/user"
	"go.uber.org/zap"
	"net/http"
)

type RegistrationHandler struct {
	service     *auth.RegistrationService
	authService *auth.Service
	logger      *zap.SugaredLogger
}

func NewRegistrationHandler(service *auth.RegistrationService, authService *auth.Service, logger *zap.SugaredLogger) *RegistrationHandler {
	return &RegistrationHandler{
		service:     service,
		authService: authService,
		logger:      logger,
	}
}

func (h *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req auth.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Info("Register request: %v", req)

	w.Header().Set("Content-Type", "application/json")

	register, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(&auth.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.authService.LoginByID(r.Context(), register.ID)
	if err != nil {
		return
	}

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&auth.RegisterResponse{ID: register.ID, Login: string(register.Login)})
}
