package schema

type Reserve struct {
	OrderID     int64  `db:"order_id"`
	WarehouseID int64  `db:"warehouse_id"`
	SKU         int64  `db:"sku"`
	Count       uint64 `db:"count"`
}
