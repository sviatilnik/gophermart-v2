package order

import orderDomain "github.com/sviatilnik/gophermart/internal/domain/order"

var (
	ErrAlreadyExists                 = orderDomain.ErrAlreadyExists
	ErrAlreadyCreatedByOtherCustomer = orderDomain.ErrAlreadyCreatedByOtherCustomer
)
