package accrual

type Service interface {
	CheckOrder(orderNumber string) (*Accrual, error)
}
