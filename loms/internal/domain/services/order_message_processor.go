package services

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	"time"

	"go.uber.org/zap"
)

type OrderMessageProcessor struct {
	publisher domain.OrderPublisher
	repo      domain.OrderRepository
}

func NewOrderMessageProcessor(repo domain.OrderRepository, publisher domain.OrderPublisher) *OrderMessageProcessor {
	return &OrderMessageProcessor{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *OrderMessageProcessor) publishMessages(ctx context.Context) error {
	rows, err := s.repo.GetOrderMessages(ctx, false)
	if err != nil {
		return err
	}

	for _, row := range rows {
		msg := domain.OrderMessage{
			OrderID:  row.OrderID,
			StatusID: row.StatusID,
		}

		isProcessed := true
		errString := ""

		err = s.publisher.PublishOrderMessage(msg)
		if err != nil {
			isProcessed = false
			errString = err.Error()
		}

		err = s.repo.UpdateOrderMessage(ctx, row.ID, isProcessed, errString)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderMessageProcessor) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			err := s.publishMessages(ctx)
			if err != nil {
				logger.Info("message processor error", zap.Error(err))
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
