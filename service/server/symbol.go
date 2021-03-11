package server

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/service"

	"github.com/vectorman1/analysis/analysis-api/common"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *SymbolsServiceServer) ReadPaged(ctx context.Context, req *symbol_service.ReadPagedSymbolRequest) (*symbol_service.ReadPagedSymbolResponse, error) {
	timeoutContext, c := context.WithTimeout(ctx, 5*time.Second)
	defer c()

	res, totalItemsCount, err := s.symbolService.GetPaged(timeoutContext, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	resp := &symbol_service.ReadPagedSymbolResponse{
		Items:      *res,
		TotalItems: uint64(totalItemsCount),
	}
	return resp, nil
}

func (s *SymbolsServiceServer) Read(ctx context.Context, req *symbol_service.ReadSymbolRequest) (*symbol_service.ReadSymbolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}

func (s *SymbolsServiceServer) Create(ctx context.Context, req *symbol_service.CreateSymbolRequest) (*symbol_service.CreateSymbolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (s *SymbolsServiceServer) Update(ctx context.Context, req *symbol_service.UpdateSymbolRequest) (*symbol_service.UpdateSymbolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *SymbolsServiceServer) Delete(ctx context.Context, req *symbol_service.DeleteSymbolRequest) (*symbol_service.DeleteSymbolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *SymbolsServiceServer) Details(ctx context.Context, req *symbol_service.SymbolDetailsRequest) (*symbol_service.SymbolDetailsResponse, error) {
	res, err := s.symbolService.Details(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}

func (s *SymbolsServiceServer) Recalculate(ctx context.Context, req *symbol_service.RecalculateSymbolRequest) (*symbol_service.RecalculateSymbolResponse, error) {
	res, err := s.symbolService.Recalculate(ctx)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}
	return res, nil
}
