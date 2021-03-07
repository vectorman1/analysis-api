package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx"

	"github.com/vectorman1/analysis/analysis-api/service"

	"github.com/vectorman1/analysis/analysis-api/service/server"

	"github.com/vectorman1/analysis/analysis-api/common"
	"github.com/vectorman1/analysis/analysis-api/db"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	grpc_server "github.com/vectorman1/analysis/analysis-api/server/grpc-server"
	rest_server "github.com/vectorman1/analysis/analysis-api/server/rest-server"
)

// RunServer runs gRPC grpc-server and HTTP gateway
func RunServer() error {
	ctx := context.Background()
	config, err := common.GetConfig()
	if err != nil {
		return err
	}

	// get configuration
	if len(config.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC grpc-server: '%s'", config.GRPCPort)
	}
	if len(config.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", config.HTTPPort)
	}

	// initialize logger-grpc
	if err := logger_grpc.Init(config.LogLevel, config.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger-grpc: %v", err)
	}

	// set up db connection pool
	dbConnPool, err := db.GetConnPool(config)
	if err != nil {
		return fmt.Errorf("failed to create conn pool: %v", err)
	}

	conn, err := dbConnPool.Acquire()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer conn.Close()

	// connect to worker rpc server
	rpcClient := common.NewRpcClient(config)
	_, err = rpcClient.Initialize()
	if err != nil {
		return fmt.Errorf("failed to connect to worker rpc server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
	}()

	s, err := initializeServices(ctx, dbConnPool, rpcClient, config, sigs)
	if err != nil {
		return err
	}

	// run HTTP gateway
	go func() {
		_ = rest_server.RunServer(ctx, config.GRPCPort, config.HTTPPort)
	}()

	return s.Run()
}

func initializeServices(ctx context.Context, dbConnPool *pgx.ConnPool, rpcClient *common.Rpc, config *common.Config, sigs chan os.Signal) (*grpc_server.GRPCServer, error) {
	//w := zerolog.NewConsoleWriter()
	//
	//symbolsQueue := common.NewRabbitClient("symbols.stream", "symbols", config.RabbitMqConn, sigs, zerolog.New(w))
	//defer symbolsQueue.Close()

	symbolRepository := db.NewSymbolRepository(dbConnPool)
	currencyRepository := db.NewCurrencyRepository(dbConnPool)
	symbolOverviewRepository := db.NewSymbolOverviewRepository(dbConnPool)
	userRepository := db.NewUserRepository(dbConnPool)
	historicalRepository := db.NewHistoricalRepository(dbConnPool)

	alphaVantageService := service.NewAlphaVantageService(config)

	symbolsService := service.NewSymbolsService(symbolRepository, symbolOverviewRepository, currencyRepository, alphaVantageService)
	userService := service.NewUserService(userRepository, config)
	historicalService := service.NewHistoricalService(historicalRepository, symbolRepository)

	symbolsServiceServer := server.NewSymbolsServiceServer(rpcClient, symbolsService)
	userServiceServer := server.NewUserServiceServer(userService)
	historicalServiceServer := server.NewHistoricalServiceServer(historicalService)

	return grpc_server.NewGRPCServer(ctx, config.GRPCPort, symbolsServiceServer, userServiceServer, historicalServiceServer), nil
}

func main() {
	if err := RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
