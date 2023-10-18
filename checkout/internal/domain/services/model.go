package services

//go:generate minimock -i TransactionManager -o ./mocks/ -s "_mock.go"

import (
	"context"
	"route256/checkout/internal/domain"
)

type Model struct {
	orderManager        OrderManager
	productManager      ProductManager
	cartRepository      CartRepository
	transactionManager  TransactionManager
	listCartWorkerCount int
}

type OrderManager interface {
	Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error)
	CreateOrder(ctx context.Context, user int64, items []domain.CartItem) (*domain.Order, error)
}

type ProductManager interface {
	GetProduct(ctx context.Context, sku uint32) (*domain.Product, error)
}

type CartRepository interface {
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	ListCart(ctx context.Context, user int64) ([]domain.CartItem, error)
	GetItem(ctx context.Context, user int64, sku uint32) (domain.CartItem, error)
	UpdateItemCount(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32) error
	ClearCart(ctx context.Context, user int64) error
}

func NewModel(orderManager OrderManager, productManager ProductManager, cartRepository CartRepository,
	transactionManager TransactionManager, listCartWorkerCount int) *Model {
	return &Model{
		orderManager:        orderManager,
		productManager:      productManager,
		cartRepository:      cartRepository,
		transactionManager:  transactionManager,
		listCartWorkerCount: listCartWorkerCount,
	}
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunSerializable(ctx context.Context, fx func(ctxTX context.Context) error) error
}
