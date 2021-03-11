package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vectorman1/analysis/analysis-api/jobs"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

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

	// set up postgres db connection pool
	dbConnPool, err := db.GetConnPool(config)
	if err != nil {
		return fmt.Errorf("failed to create entities conn pool: %v", err)
	}

	conn, err := dbConnPool.Acquire()
	if err != nil {
		return fmt.Errorf("failed to open entities database: %v", err)
	}
	defer conn.Close()

	// set up mongodb connection pool
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDbConnString))
	if err != nil {
		return fmt.Errorf("failed to create documents client: %v", err)
	}

	tctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	err = client.Connect(tctx)
	if err != nil {
		return fmt.Errorf("failed to create documents conn pool: %v", err)
	}
	defer client.Disconnect(tctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
	}()

	s, err := initializeServices(ctx, dbConnPool, client.Database(common.MONGO_DB_DATABASE), config)
	if err != nil {
		return err
	}

	// run HTTP gateway
	go func() {
		_ = rest_server.RunServer(ctx, config)
	}()

	return s.Run()
}

func initializeServices(ctx context.Context, pgConnPool *pgx.ConnPool, mongoDatabase *mongo.Database, config *common.Config) (*grpc_server.GRPCServer, error) {
	historicalRepository := db.NewHistoricalRepository(pgConnPool, mongoDatabase)
	symbolOverviewRepository := db.NewSymbolOverviewRepository(pgConnPool, mongoDatabase)

	symbolRepository := db.NewSymbolRepository(pgConnPool)
	userRepository := db.NewUserRepository(pgConnPool)

	externalSymbolService := service.NewExternalSymbolService()
	alphaVantageService := service.NewAlphaVantageService(config)

	symbolsService := service.NewSymbolsService(symbolRepository, symbolOverviewRepository, alphaVantageService, externalSymbolService)
	userService := service.NewUserService(userRepository, config)
	historicalService := service.NewHistoricalService(historicalRepository, symbolRepository, symbolOverviewRepository)

	symbolsServiceServer := server.NewSymbolsServiceServer(symbolsService)
	userServiceServer := server.NewUserServiceServer(userService)
	historicalServiceServer := server.NewHistoricalServiceServer(historicalService)

	err := jobs.ScheduleJobs(symbolsService)
	if err != nil {
		return nil, err
	}

	return grpc_server.NewGRPCServer(ctx, config.GRPCPort, symbolsServiceServer, userServiceServer, historicalServiceServer), nil
}

func main() {
	if err := RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
