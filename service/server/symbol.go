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
	symbolService *service.SymbolsService
	symbol_service.UnimplementedSymbolServiceServer
}

func NewSymbolsServiceServer(
	symbolsService *service.SymbolsService) *SymbolsServiceServer {
	return &SymbolsServiceServer{
		symbolService: symbolsService,
	}
}

func (s *SymbolsServiceServer) ReadPaged(ctx context.Context, req *symbol_service.ReadPagedSymbolRequest) (*symbol_service.ReadPagedSymbolResponse, error) {
	if req.Filter == nil {
		return nil, status.Errorf(codes.InvalidArgument, "provide filter")
	}
	if req.Filter.Order == "" {
		return nil, status.Error(codes.InvalidArgument, "provide order argument")
	}
	timeoutContext, c := context.WithTimeout(ctx, 5*time.Second)
	defer c()

	res, totalItemsCount, err := s.symbolService.GetPaged(timeoutContext, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
	userInfo := ctx.Value("user_info")
	if userInfo == nil {
		return nil, status.Error(codes.Unauthenticated, "provide user token")
	}

	res, err := s.symbolService.Details(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, st.Err()
		}
		return nil, err
	}

	return res, nil
}

func (s *SymbolsServiceServer) Recalculate(ctx context.Context, req *symbol_service.RecalculateSymbolRequest) (*symbol_service.RecalculateSymbolResponse, error) {
	userInfo := ctx.Value("user_info")
	if userInfo == nil {
		return nil, status.Error(codes.Unauthenticated, "provide user token")
	}

	res, err := s.symbolService.Recalculate(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
