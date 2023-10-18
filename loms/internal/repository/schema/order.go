package schema

type OrderStatus int32

const (
	OrderStatusNew OrderStatus = iota + 1
	OrderStatusAwaitingPayment
	OrderStatusFailed
	OrderStatusPayed
	OrderStatusCancelled
)

type Order struct {
	ID       int64       `db:"id"`
	UserID   int64       `db:"user_id"`
	StatusID OrderStatus `db:"status_id"`
}
