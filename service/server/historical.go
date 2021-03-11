package server

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"
	"github.com/vectorman1/analysis/analysis-api/service"
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
	res, err := s.historicalService.GetBySymbolUuid(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}
