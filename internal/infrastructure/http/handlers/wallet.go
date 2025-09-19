package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/gophermart/internal/application/wallet"
	"github.com/sviatilnik/gophermart/internal/infrastructure/http/middleware"
	"net/http"
)

type WalletHandler struct {
	service *wallet.Service
}

func NewWalletHandler(service *wallet.Service) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) Balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.Context().Value(middleware.RequestUserID).(string)

	wlt, err := h.service.Get(r.Context(), customerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&wlt)
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.Context().Value(middleware.RequestUserID).(string)

	var withdrawal wallet.CreateWithdrawal
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	err = h.service.Withdraw(r.Context(), customerID, withdrawal.OrderNumber, withdrawal.Amount)
	if err != nil {
		if errors.Is(err, wallet.ErrNotEnoughFunds) {
			w.WriteHeader(http.StatusPaymentRequired)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WalletHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.Context().Value(middleware.RequestUserID).(string)

	withdraws, err := h.service.GetWithdraws(r.Context(), customerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	if len(withdraws) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(withdraws)
}
