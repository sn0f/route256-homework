package main

import (
	"context"
	"fmt"
	"sync"

	"route256/libs/kafka"
	"route256/libs/logger"
	"route256/notifications/internal/config"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func main() {
	keepRunning := true
	err := logger.Init(false)
	if err != nil {
		panic(fmt.Errorf("error creating logger: %v", err))
	}
	logger.Info("starting a new sarama consumer")

	err = config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.MaxVersion
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	handler := func(message *sarama.ConsumerMessage) {
		logger.Info("message received",
			zap.String("value", string(message.Value)),
			zap.Time("timestamp", message.Timestamp),
			zap.String("topic", message.Topic))
	}
	consumer := kafka.NewConsumerGroup(handler)

	ctx, cancel := context.WithCancel(context.Background())
	group, err := sarama.NewConsumerGroup(config.ConfigData.Brokers, config.ConfigData.Group, saramaConfig)
	if err != nil {
		panic(fmt.Errorf("error creating consumer group client: %v", err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := group.Consume(ctx, []string{config.ConfigData.Topic}, &consumer); err != nil {
				panic(fmt.Errorf("error from consumer: %v", err))
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	logger.Info("sarama consumer up and running!...")

	for keepRunning {
		<-ctx.Done()
		logger.Info("terminating: context cancelled")
		keepRunning = false
	}

	cancel()
	wg.Wait()

	if err = group.Close(); err != nil {
		panic(fmt.Errorf("error closing sarama client: %v", err))
	}
}
