package grpc_server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"
	"github.com/vectorman1/analysis/analysis-api/generated/user_service"

	instrument_present "github.com/vectorman1/analysis/analysis-api/domain/instrument/present"
	user_present "github.com/vectorman1/analysis/analysis-api/domain/user/present"

	"github.com/vectorman1/analysis/analysis-api/middleware"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Context             context.Context
	Port                string
	symbolServiceServer *instrument_present.InstrumentServiceServer
	userServiceServer   *user_present.UserServiceServer
}

func NewGRPCServer(
	ctx context.Context,
	port string,
	instrumentServiceServer *instrument_present.InstrumentServiceServer,
	userServiceServer *user_present.UserServiceServer) *GRPCServer {
	return &GRPCServer{
		Context:             ctx,
		Port:                port,
		symbolServiceServer: instrumentServiceServer,
		userServiceServer:   userServiceServer,
	}
}

// Run runs gRPC service to publish our services
func (s *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}

	// add middleware
	opts := middleware.LoadMiddleware(logger_grpc.Log, nil)

	server := grpc.NewServer(opts...)

	// register services
	instrument_service.RegisterInstrumentServiceServer(server, s.symbolServiceServer)
	user_service.RegisterUserServiceServer(server, s.userServiceServer)

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
