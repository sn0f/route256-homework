package postgres

import (
	"context"
	"route256/libs/postgres"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type orderRepo struct {
	postgres.QueryEngineProvider
}

func NewOrderRepo(provider postgres.QueryEngineProvider) *orderRepo {
	return &orderRepo{
		QueryEngineProvider: provider,
	}
}

var (
	tableOrders        = "orders"
	tableOrderItems    = "order_items"
	tableStocks        = "stocks"
	tableReserves      = "reserves"
	tableOrderMessages = "order_messages"
	pgBuilder          = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func (r *orderRepo) GetStocks(ctx context.Context, sku uint32) ([]domain.StockItem, error) {
	sql, args, err := pgBuilder.Select("warehouse_id", "count").
		From(tableStocks).
		Where(sq.Eq{"sku": sku}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.Stock
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, db, &items, sql, args...); err != nil {
		return nil, err
	}

	return schema.BindSchemaStocksToModelStockItems(items), nil
}

func (r *orderRepo) UpdateStocks(ctx context.Context, warehouseID int64, sku uint32, count uint64) error {
	sql, args, err := pgBuilder.Update(tableStocks).Set("count", count).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku}).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if _, err = db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) SubtractStocks(ctx context.Context, warehouseID int64, sku uint32, count uint64) error {
	sql, args, err := pgBuilder.Update(tableStocks).Set("count", sq.Expr("count-?", count)).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku}).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if _, err = db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) CreateReserve(ctx context.Context, reserve domain.Reserve) error {
	sql, args, err := pgBuilder.Insert(tableReserves).
		Columns("order_id", "warehouse_id", "sku", "count").
		Values(reserve.OrderID, reserve.WarehouseID, reserve.SKU, reserve.Count).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) GetReservedCount(ctx context.Context, warehouseID int64, sku uint32) (uint64, error) {
	sql, args, err := pgBuilder.Select("coalesce(SUM(reserves.count), 0) as count").From(tableReserves).
		Join("orders on reserves.order_id = orders.id").
		Where(sq.Eq{"orders.status_id": schema.OrderStatusAwaitingPayment}).
		Where(sq.Eq{"reserves.warehouse_id": warehouseID}).
		Where(sq.Eq{"reserves.sku": sku}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *orderRepo) GetReserves(ctx context.Context, orderID int64) ([]domain.Reserve, error) {
	sql, args, err := pgBuilder.Select("warehouse_id", "sku", "count").
		From(tableReserves).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.Reserve
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = pgxscan.Select(ctx, db, &items, sql, args...)
	if err != nil {
		return nil, err
	}

	return schema.BindSchemaReservesToModelReserves(items), nil
}

func (r *orderRepo) DeleteReserves(ctx context.Context, orderID int64) error {
	sql, args, err := pgBuilder.Delete(tableReserves).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) CreateOrder(ctx context.Context, user int64) (int64, error) {
	sql, args, err := pgBuilder.Insert(tableOrders).
		Columns("user_id", "status_id").
		Values(user, schema.OrderStatusNew).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var order schema.Order
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = pgxscan.Get(ctx, db, &order, sql, args...)
	if err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (r *orderRepo) InsertOrderItems(ctx context.Context, orderID int64, items []domain.OrderItem) error {
	query := pgBuilder.Insert(tableOrderItems).
		Columns("order_id", "sku", "count")
	for _, item := range items {
		query = query.Values(orderID, item.SKU, item.Count)
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) GetOrder(ctx context.Context, orderID int64) (domain.Order, error) {
	sql, args, err := pgBuilder.Select("id", "user_id", "status_id").
		From(tableOrders).
		Where(sq.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return domain.Order{}, err
	}

	var order schema.Order
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = pgxscan.Get(ctx, db, &order, sql, args...)
	if err != nil {
		return domain.Order{}, err
	}

	return schema.BindSchemaOrderToModelOrder(order), nil
}

func (r *orderRepo) GetOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	sql, args, err := pgBuilder.Select("sku", "count").
		From(tableOrderItems).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.OrderItem
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, db, &items, sql, args...); err != nil {
		return nil, err
	}

	return schema.BindSchemaOrderItemsToModelOrderItems(items), nil
}

func (r *orderRepo) UpdateOrder(ctx context.Context, orderID int64, statusID domain.OrderStatus) error {
	sql, args, err := pgBuilder.Update(tableOrders).Set("status_id", schema.OrderStatus(statusID)).
		Where(sq.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) CreateOrderMessage(ctx context.Context, orderID int64, statusID domain.OrderStatus) (int64, error) {
	sql, args, err := pgBuilder.Insert(tableOrderMessages).
		Columns("order_id", "status_id").
		Values(orderID, schema.OrderStatus(statusID)).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var msg schema.OrderMessage
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = pgxscan.Get(ctx, db, &msg, sql, args...)
	if err != nil {
		return 0, err
	}

	return msg.ID, nil
}

func (r *orderRepo) UpdateOrderMessage(ctx context.Context, id int64, isProcessed bool, errString string) error {
	sql, args, err := pgBuilder.Update(tableOrderMessages).
		Set("is_processed", isProcessed).
		Set("error", errString).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) GetOrderMessages(ctx context.Context, isProcessed bool) ([]domain.OrderMessage, error) {
	sql, args, err := pgBuilder.Select("id", "order_id", "status_id").
		From(tableOrderMessages).
		Where(sq.Eq{"is_processed": isProcessed}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.OrderMessage
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, db, &items, sql, args...); err != nil {
		return nil, err
	}

	return schema.BindSchemaOrderMessagesToModelOrderMessagess(items), nil
}
