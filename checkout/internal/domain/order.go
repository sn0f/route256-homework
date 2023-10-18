package domain

type Order struct {
	OrderID int64
}

func NewOrder(orderID int64) *Order {
	return &Order{OrderID: orderID}
}
