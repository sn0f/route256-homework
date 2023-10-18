package services

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Сервис аннулирования заказов по таймауту
type CancelOrderService struct {
	// Время жизни заявки
	orderLifeTime time.Duration
	// Мапа для хранения таймеров
	mapCancels map[int64]*time.Timer
	// Для доступа в мапу
	mutex *sync.RWMutex
}

func NewOrderCanceller(orderLifeTime time.Duration) domain.OrderCanceller {
	return &CancelOrderService{
		orderLifeTime: orderLifeTime,
		mapCancels:    make(map[int64]*time.Timer),
		mutex:         &sync.RWMutex{},
	}
}

// Запуск задачи аннулирования заказа по заданному таймауту
func (s CancelOrderService) AddCancelOrderTask(orderID int64, m *domain.Model) {
	ctx := context.Background()

	timer := time.AfterFunc(s.orderLifeTime, func() {
		logger.Info("cancelling order after timeout", zap.Int64("order_id", orderID))
		s.RemoveCancelOrderTask(orderID)

		err := m.CancelOrder(ctx, orderID)
		if err != nil {
			logger.Info("cancelling order by timeout failed", zap.Error(err))
		}
		logger.Info("order cancelled after timeout", zap.Int64("order_id", orderID))
	})

	// Сохраняем таймер в мапе
	s.mutex.Lock()
	s.mapCancels[orderID] = timer
	s.mutex.Unlock()
}

// Отменяем контекст и удаляем таймер из мапы
func (s CancelOrderService) RemoveCancelOrderTask(orderID int64) {
	s.mutex.RLock()
	timer, ok := s.mapCancels[orderID]
	if !ok {
		s.mutex.RUnlock()
		return
	}
	s.mutex.RUnlock()
	timer.Stop()

	s.mutex.Lock()
	delete(s.mapCancels, orderID)
	s.mutex.Unlock()
	logger.Info("task for cancelling order successfully removed", zap.Int64("order_id", orderID))
}
