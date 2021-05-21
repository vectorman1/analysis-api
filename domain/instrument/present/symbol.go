package present

import (
	"context"
	"time"

	service2 "github.com/vectorman1/analysis/analysis-api/domain/instrument/service"

	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"

	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/jobs"

	"github.com/vectorman1/analysis/analysis-api/common"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type SymbolsServiceServer struct {
	rabbitClient  *common.RabbitClient
	symbolService *service2.SymbolService
	symbol_service.UnimplementedSymbolServiceServer
}

func NewSymbolServiceServer(
	symbolsService *service2.SymbolService) *SymbolsServiceServer {
	return &SymbolsServiceServer{
		symbolService: symbolsService,
	}
}

func (s *SymbolsServiceServer) Get(ctx context.Context, req *symbol_service.SymbolRequest) (*proto_models.Symbol, error) {
	return s.symbolService.Get(ctx, req.Uuid)
}

func (s *SymbolsServiceServer) GetPaged(ctx context.Context, req *symbol_service.GetPagedRequest) (*symbol_service.GetPagedResponse, error) {
	timeoutContext, c := context.WithTimeout(ctx, 5*time.Second)
	defer c()

	res, totalItemsCount, err := s.symbolService.GetPaged(timeoutContext, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	resp := &symbol_service.GetPagedResponse{
		Items:      *res,
		TotalItems: uint64(totalItemsCount),
	}
	return resp, nil
}

func (s *SymbolsServiceServer) Overview(ctx context.Context, req *symbol_service.SymbolOverviewRequest) (*symbol_service.SymbolOverview, error) {
	res, err := s.symbolService.Overview(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}

func (s *SymbolsServiceServer) UpdateAllJob(ctx context.Context, req *symbol_service.StartUpdateJobRequest) (*symbol_service.StartUpdateJobResponse, error) {
	jobrunner.Now(jobs.NewSymbolUpdateJob(s.symbolService))

	return &symbol_service.StartUpdateJobResponse{}, nil
}

func (s *SymbolsServiceServer) UpdateAll(ctx context.Context, req *symbol_service.StartUpdateJobRequest) (*symbol_service.UpdateAllResponse, error) {
	res, err := s.symbolService.UpdateAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
