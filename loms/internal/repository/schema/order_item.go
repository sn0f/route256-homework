package schema

type OrderItem struct {
	ID      int64  `db:"id"`
	OrderID int64  `db:"order_id"`
	SKU     int64  `db:"sku"`
	Count   uint64 `db:"count"`
}
