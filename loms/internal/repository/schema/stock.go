package schema

type Stock struct {
	WarehouseID int64  `db:"warehouse_id"`
	SKU         int64  `db:"sku"`
	Count       uint64 `db:"count"`
}
