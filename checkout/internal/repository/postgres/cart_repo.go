package postgres

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"
	"route256/libs/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type cartRepo struct {
	postgres.QueryEngineProvider
}

func NewCartRepo(provider postgres.QueryEngineProvider) *cartRepo {
	return &cartRepo{
		QueryEngineProvider: provider,
	}
}

var (
	tableCarts = "cart_items"
	pgBuilder  = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func (r *cartRepo) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	sql, args, err := pgBuilder.Insert(tableCarts).
		Columns("user_id", "sku", "count").
		Values(user, sku, count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = cart_items.count + EXCLUDED.count").
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

func (r *cartRepo) DeleteFromCart(ctx context.Context, user int64, sku uint32) error {
	sql, args, err := pgBuilder.Delete(tableCarts).
		Where(sq.Eq{"user_id": user}).
		Where(sq.Eq{"sku": sku}).
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

func (r *cartRepo) ClearCart(ctx context.Context, user int64) error {
	sql, args, err := pgBuilder.Delete(tableCarts).
		Where(sq.Eq{"user_id": user}).
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

func (r *cartRepo) UpdateItemCount(ctx context.Context, user int64, sku uint32, count uint16) error {
	sql, args, err := pgBuilder.Update(tableCarts).Set("count", count).
		Where(sq.Eq{"user_id": user}).
		Where(sq.Eq{"sku": sku}).
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

func (r *cartRepo) GetItem(ctx context.Context, user int64, sku uint32) (domain.CartItem, error) {
	sql, args, err := pgBuilder.Select("sku", "count").From(tableCarts).
		Where(sq.Eq{"user_id": user}).
		Where(sq.Eq{"sku": sku}).
		ToSql()
	if err != nil {
		return domain.CartItem{}, err
	}

	var item domain.CartItem
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	err = pgxscan.Get(ctx, db, &item, sql, args...)
	if err != nil {
		return domain.CartItem{}, err
	}

	return item, nil
}

func (r *cartRepo) ListCart(ctx context.Context, user int64) ([]domain.CartItem, error) {
	sql, args, err := pgBuilder.Select("sku", "count").From(tableCarts).
		Where(sq.Eq{"user_id": user}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.CartItem
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, db, &items, sql, args...); err != nil {
		return nil, err
	}

	return schema.BindSchemaCartItemsToModelCartItems(items), nil
}
