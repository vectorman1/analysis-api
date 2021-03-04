package grpc_server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"

	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
	"github.com/vectorman1/analysis/analysis-api/service/server"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/middleware"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Context                 context.Context
	Port                    string
	symbolsServiceServer    symbol_service.SymbolServiceServer
	userServiceServer       user_service.UserServiceServer
	historicalServiceServer historical_service.HistoricalServiceServer
}

func NewGRPCServer(
	ctx context.Context,
	port string,
	symbolsServiceServer symbol_service.SymbolServiceServer,
	userServiceServer *server.UserServiceServer,
	historicalServiceServer *server.HistoricalServiceServer) *GRPCServer {
	return &GRPCServer{
		Context:                 ctx,
		Port:                    port,
		symbolsServiceServer:    symbolsServiceServer,
		userServiceServer:       userServiceServer,
		historicalServiceServer: historicalServiceServer,
	}
}

// RunServer runs gRPC service to publish our services
func (s *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}

	// add middleware
	opts := middleware.LoadMiddleware(logger_grpc.Log, nil)

	server := grpc.NewServer(opts...)

	// register services
	symbol_service.RegisterSymbolServiceServer(server, s.symbolsServiceServer)
	user_service.RegisterUserServiceServer(server, s.userServiceServer)
	historical_service.RegisterHistoricalServiceServer(server, s.historicalServiceServer)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC grpc-server...")

			server.GracefulStop()

			<-s.Context.Done()
		}
	}()

	// start gRPC grpc-server
	log.Println("starting gRPC grpc-server...")
	return server.Serve(listen)
}
