package server

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"
	"github.com/vectorman1/analysis/analysis-api/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HistoricalServiceServer struct {
	historicalService *service.HistoricalService
	historical_service.UnimplementedHistoricalServiceServer
}

func NewHistoricalServiceServer(historicalService *service.HistoricalService) *HistoricalServiceServer {
	return &HistoricalServiceServer{
		historicalService: historicalService,
	}
}

func (s *HistoricalServiceServer) GetBySymbolUuid(ctx context.Context, req *historical_service.GetBySymbolUuidRequest) (*historical_service.GetBySymbolUuidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, req.SymbolUuid)
}
