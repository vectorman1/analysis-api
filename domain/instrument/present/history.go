package present

import (
	"context"

	service2 "github.com/vectorman1/analysis/analysis-api/domain/instrument/service"

	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/jobs"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/vectorman1/analysis/analysis-api/common"
)

type HistoryServiceServer struct {
	historyService *service2.HistoryService
	history_service.UnimplementedHistoryServiceServer
}

func NewHistoryServiceServer(historicalService *service2.HistoryService) *HistoryServiceServer {
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

func (s *HistoryServiceServer) StartUpdateJob(ctx context.Context, req *history_service.StartUpdateJobRequest) (*history_service.StartUpdateJobResponse, error) {
	jobrunner.Now(jobs.NewHistoryUpdateJob(s.historyService))

	return &history_service.StartUpdateJobResponse{}, nil
}

func (s *HistoryServiceServer) GetChartBySymbolUuid(ctx context.Context, req *history_service.GetChartBySymbolUuidRequest) (*history_service.GetChartBySymbolUuidResponse, error) {
	res, err := s.historyService.GetChartBySymbolUuid(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
