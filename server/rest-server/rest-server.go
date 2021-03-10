package rest_server

import (
	"context"
	"strings"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"

	"github.com/vectorman1/analysis/analysis-api/generated/user_service"

	logger_rest "github.com/vectorman1/analysis/analysis-api/middleware/logger-rest"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	tracer_rest "github.com/vectorman1/analysis/analysis-api/middleware/tracer-rest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, config *common.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := symbol_service.RegisterSymbolServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:"+config.GRPCPort, opts); err != nil {
		logger_grpc.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}
	if err := user_service.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:"+config.GRPCPort, opts); err != nil {
		logger_grpc.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}
	if err := historical_service.RegisterHistoricalServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:"+config.GRPCPort, opts); err != nil {
		logger_grpc.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + config.HTTPPort,
		Handler: tracer_rest.AddRequestID(logger_rest.AddLogger(logger_grpc.Log, allowCORS(mux, config.AllowedOrigin))),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Println("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(common.GetAllowedHeaders(), ","))
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(common.GetAllowedMethods(), ","))
	return
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
