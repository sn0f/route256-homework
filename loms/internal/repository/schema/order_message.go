package schema

type OrderMessage struct {
	ID          int64       `db:"id"`
	OrderID     int64       `db:"order_id"`
	StatusID    OrderStatus `db:"status_id"`
	IsProcessed bool        `db:"is_processed"`
	Error       string      `db:"error"`
}
