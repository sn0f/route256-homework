package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain/services"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/internal/server/checkout"
	api "route256/checkout/pkg/checkout/v1"
	clientInterceptors "route256/libs/client/interceptors"
	"route256/libs/logger"
	libs "route256/libs/postgres"
	"route256/libs/server/interceptors"
	"route256/libs/server/rate"
	"route256/libs/server/tracing"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const port = "8080"

func main() {
	err := logger.Init(false)
	if err != nil {
		panic(fmt.Errorf("error creating logger: %v", err))
	}
	logger.Info("starting checkout service")

	tracing.Init(logger.Get(), "checkout")

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

	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			clientInterceptors.MetricsInterceptor,
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		logger.Fatal("failed to connect to loms server", zap.Error(err))
	}
	defer lomsConn.Close()

	productConn, err := grpc.Dial(config.ConfigData.Services.Products,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			clientInterceptors.MetricsInterceptor,
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		logger.Fatal("failed to connect to product server", zap.Error(err))
	}
	defer productConn.Close()

	txManager := libs.NewTransactionManager(pool)
	cartRepo := postgres.NewCartRepo(txManager)
	rateLimiter := rate.NewLimiter(config.ConfigData.ProductsRpsLimit, config.ConfigData.ProductsBurstLimit)
	productClient := product.New(productConn, config.ConfigData.Token, rateLimiter,
		time.Duration(config.ConfigData.ProductsCacheTTL), config.ConfigData.ProductsCacheSize)
	lomsClient := loms.New(lomsConn)
	model := services.NewModel(lomsClient, productClient, cartRepo, txManager, config.ConfigData.ListCartWorkerCount)

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

	api.RegisterCheckoutServiceServer(s, checkout.NewCheckoutServer(model))

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
