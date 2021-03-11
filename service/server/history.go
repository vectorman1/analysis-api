package server

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/vectorman1/analysis/analysis-api/service"
)

type HistoryServiceServer struct {
	historyService *service.HistoryService
	history_service.UnimplementedHistoryServiceServer
}

func NewHistoryServiceServer(historicalService *service.HistoryService) *HistoryServiceServer {
	return &HistoryServiceServer{
		historyService: historicalService,
	}
}

func (s *HistoryServiceServer) GetBySymbolUuid(ctx context.Context, req *history_service.GetBySymbolUuidRequest) (*history_service.GetBySymbolUuidResponse, error) {
	res, err := s.historyService.GetSymbolHistory(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}
