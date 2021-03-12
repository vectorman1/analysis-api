package server

import (
	"context"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/jobs"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/service"

	"github.com/vectorman1/analysis/analysis-api/common"
)

type SymbolsServiceServer struct {
	rabbitClient  *common.RabbitClient
	symbolService *service.SymbolService
	symbol_service.UnimplementedSymbolServiceServer
}

func NewSymbolServiceServer(
	symbolsService *service.SymbolService) *SymbolsServiceServer {
	return &SymbolsServiceServer{
		symbolService: symbolsService,
	}
}

func (s *SymbolsServiceServer) ReadPaged(ctx context.Context, req *symbol_service.ReadPagedRequest) (*symbol_service.ReadPagedResponse, error) {
	timeoutContext, c := context.WithTimeout(ctx, 5*time.Second)
	defer c()

	res, totalItemsCount, err := s.symbolService.GetPaged(timeoutContext, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	resp := &symbol_service.ReadPagedResponse{
		Items:      *res,
		TotalItems: uint64(totalItemsCount),
	}
	return resp, nil
}

func (s *SymbolsServiceServer) Details(ctx context.Context, req *symbol_service.SymbolDetailsRequest) (*symbol_service.SymbolDetailsResponse, error) {
	res, err := s.symbolService.Details(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}

func (s *SymbolsServiceServer) StartUpdateJob(ctx context.Context, req *symbol_service.StartUpdateJobRequest) (*symbol_service.StartUpdateJobResponse, error) {
	jobrunner.Now(jobs.NewSymbolUpdateJob(s.symbolService))

	return &symbol_service.StartUpdateJobResponse{}, nil
}
