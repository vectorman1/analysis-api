package grpc_server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
	"github.com/vectorman1/analysis/analysis-api/service/server"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/middleware"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Context              context.Context
	Port                 string
	symbolServiceServer  symbol_service.SymbolServiceServer
	userServiceServer    user_service.UserServiceServer
	historyServiceServer history_service.HistoryServiceServer
}

func NewGRPCServer(
	ctx context.Context,
	port string,
	symbolServiceServer symbol_service.SymbolServiceServer,
	userServiceServer *server.UserServiceServer,
	historyServiceServer *server.HistoryServiceServer) *GRPCServer {
	return &GRPCServer{
		Context:              ctx,
		Port:                 port,
		symbolServiceServer:  symbolServiceServer,
		userServiceServer:    userServiceServer,
		historyServiceServer: historyServiceServer,
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
	symbol_service.RegisterSymbolServiceServer(server, s.symbolServiceServer)
	user_service.RegisterUserServiceServer(server, s.userServiceServer)
	history_service.RegisterHistoryServiceServer(server, s.historyServiceServer)

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
