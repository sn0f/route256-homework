package domain

import (
	"context"
)

type Model struct {
	orderRepository    OrderRepository
	transactionManager TransactionManager
	orderCanceller     OrderCanceller
}

type OrderRepository interface {
	GetStocks(ctx context.Context, sku uint32) ([]StockItem, error)
	UpdateStocks(ctx context.Context, warehouseID int64, sku uint32, count uint64) error
	SubtractStocks(ctx context.Context, warehouseID int64, sku uint32, count uint64) error

	CreateReserve(ctx context.Context, reserve Reserve) error
	GetReservedCount(ctx context.Context, warehouseID int64, sku uint32) (uint64, error)
	GetReserves(ctx context.Context, orderID int64) ([]Reserve, error)
	DeleteReserves(ctx context.Context, orderID int64) error

	CreateOrder(ctx context.Context, user int64) (int64, error)
	InsertOrderItems(ctx context.Context, orderID int64, items []OrderItem) error
	GetOrder(ctx context.Context, orderID int64) (Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]OrderItem, error)
	UpdateOrder(ctx context.Context, orderID int64, status OrderStatus) error

	CreateOrderMessage(ctx context.Context, orderID int64, statusID OrderStatus) (int64, error)
	UpdateOrderMessage(ctx context.Context, id int64, isProcessed bool, errString string) error
	GetOrderMessages(ctx context.Context, isProcessed bool) ([]OrderMessage, error)
}

func New(orderRepository OrderRepository, transactionManager TransactionManager, orderCanceller OrderCanceller) *Model {
	return &Model{
		orderRepository:    orderRepository,
		transactionManager: transactionManager,
		orderCanceller:     orderCanceller,
	}
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunSerializable(ctx context.Context, fx func(ctxTX context.Context) error) error
}

type OrderCanceller interface {
	// Запуск задачи аннулирования заказа по заданному таймауту
	AddCancelOrderTask(orderID int64, m *Model)
	// Удалить задачу аннулирования
	RemoveCancelOrderTask(orderID int64)
}

type OrderPublisher interface {
	PublishOrderMessage(OrderMessage) error
}

type Order struct {
	StatusID OrderStatus
	User     int64
	Items    []OrderItem
}

type OrderStatus int32

const (
	OrderStatusNew OrderStatus = iota + 1
	OrderStatusAwaitingPayment
	OrderStatusFailed
	OrderStatusPayed
	OrderStatusCancelled
)

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type StockItem struct {
	WarehouseID int64
	Count       uint64
}

type Reserve struct {
	OrderID     int64
	WarehouseID int64
	SKU         uint32
	Count       uint64
}

type OrderMessage struct {
	ID       int64
	OrderID  int64
	StatusID OrderStatus
}
