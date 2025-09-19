package accrual

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/gophermart/internal/domain/accrual"
	"net/http"
	"time"
)

type checkOrderResponse struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Amount      float64 `json:"accrual"`
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: time.Second * 5},
	}
}

func (c *HTTPClient) CheckOrder(orderNumber string) (*accrual.Accrual, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/orders/"+orderNumber, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	jsonResponse := &checkOrderResponse{}
	err = json.NewDecoder(response.Body).Decode(jsonResponse)
	if err != nil {
		return nil, err
	}

	return &accrual.Accrual{
		OrderNumber: jsonResponse.OrderNumber,
		State:       accrual.State(jsonResponse.Status),
		Amount:      jsonResponse.Amount,
	}, nil
}
