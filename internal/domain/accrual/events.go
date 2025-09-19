package accrual

type CreatedEvent struct {
	OrderNumber string
	Amount      float64
	Status      string
	CustomerID  string
}

func (c *CreatedEvent) GetName() string {
	return "accrual.created"
}
