package messaging

import (
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	api "route256/loms/pkg/order/v1"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type orderPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

type Handler func(id string)

func NewOrderPublisher(producer sarama.SyncProducer, topic string) *orderPublisher {
	return &orderPublisher{
		producer: producer,
		topic:    topic,
	}
}

func (p orderPublisher) PublishOrderMessage(m domain.OrderMessage) error {
	data := &api.OrderMessage{
		OrderId: m.OrderID,
		Status:  api.OrderStatus(m.StatusID),
	}

	bytes, err := protojson.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(data.OrderId)),
		Value:     sarama.ByteEncoder(bytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logger.Info("published message", zap.Int64("order_id", data.OrderId), zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
