package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/gophermart/internal/application/order"
	order2 "github.com/sviatilnik/gophermart/internal/domain/order"
	"github.com/sviatilnik/gophermart/internal/infrastructure/http/middleware"
	"io"
	"net/http"
)

type OrderHandler struct {
	service *order.Service
}

func NewOrderHandler(service *order.Service) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orderNumber, err := io.ReadAll(r.Body)
	if err != nil || string(orderNumber) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := order.CreateOrderDTO{
		Number:     string(orderNumber),
		CustomerID: r.Context().Value(middleware.RequestUserID).(string),
	}

	newOrder, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, order.ErrAlreadyExists) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, order2.ErrOrderNumberNotValid) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, order.ErrAlreadyCreatedByOtherCustomer) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newOrder)
}

func (h *OrderHandler) GetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	customerID := r.Context().Value(middleware.RequestUserID).(string)

	orders, err := h.service.GetOrders(r.Context(), customerID, 20, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
