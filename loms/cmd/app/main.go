package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"route256/libs/kafka"
	"route256/libs/logger"
	libs "route256/libs/postgres"
	"route256/libs/server/interceptors"
	"route256/libs/server/tracing"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/domain/services"
	"route256/loms/internal/messaging"
	"route256/loms/internal/repository/postgres"
	"time"

	"route256/loms/internal/server/loms"
	api "route256/loms/pkg/loms/v1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = "8081"

func main() {
	err := logger.Init(false)
	if err != nil {
		panic(fmt.Errorf("error creating logger: %v", err))
	}
	logger.Info("starting loms service")

	tracing.Init(logger.Get(), "loms")

	err = config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	connString := os.Getenv("DATABASE_URL")

	ctx := context.Background()
	pool, err := libs.CreateDBPool(ctx, connString, 60)
	if err != nil {
		logger.Fatal("failed to create database pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to acquire connection from database pool", zap.Error(err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		interceptors.LoggingInterceptor,
		interceptors.MetricsInterceptor,
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	)))
	reflection.Register(s)

	txManager := libs.NewTransactionManager(pool)
	orderRepo := postgres.NewOrderRepo(txManager)
	orderCanceller := services.NewOrderCanceller(time.Minute * 10)
	model := domain.New(orderRepo, txManager, orderCanceller)

	producer, err := kafka.NewSyncProducer(config.ConfigData.Brokers)
	if err != nil {
		logger.Fatal("failed to create sarama producer", zap.Error(err))
	}

	topic := config.ConfigData.Topic
	logger.Info("running message processor", zap.String("topic", topic))
	publisher := messaging.NewOrderPublisher(producer, topic)
	processor := services.NewOrderMessageProcessor(orderRepo, publisher)
	go processor.Run(ctx)

	api.RegisterLomsServiceServer(s, loms.NewLomsServer(model))

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		logger.Info("http server listening", zap.String("port", config.ConfigData.MetricsPort))
		err = http.ListenAndServe(fmt.Sprintf(":%v", config.ConfigData.MetricsPort), nil)
		logger.Fatal("cannot listen http", zap.Error(err))
	}()

	logger.Info("grpc server listening", zap.String("port", port))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
