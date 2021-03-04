package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/worker_symbol_service"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/service"

	"github.com/vectorman1/analysis/analysis-api/common"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SymbolsServiceServer struct {
	rpcClient     *common.Rpc
	rabbitClient  *common.RabbitClient
	symbolService *service.SymbolsService
	symbol_service.UnimplementedSymbolServiceServer
}

func NewSymbolsServiceServer(
	rpcClient *common.Rpc,
	rabbitClient *common.RabbitClient,
	symbolsService *service.SymbolsService) *SymbolsServiceServer {
	return &SymbolsServiceServer{
		rpcClient:     rpcClient,
		rabbitClient:  rabbitClient,
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
	res, err := s.symbolService.Details(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *SymbolsServiceServer) Recalculate(ctx context.Context, req *symbol_service.RecalculateSymbolRequest) (*symbol_service.RecalculateSymbolResponse, error) {
	grpcClientContext, c1 := context.WithTimeout(ctx, 60*time.Second)
	defer c1()

	s.rabbitClient.Push([]byte("1234"))
	client := worker_symbol_service.NewWorkerSymbolServiceClient(s.rpcClient.Connection)
	stream, err := client.RecalculateSymbols(grpcClientContext)

	if err != nil {
		return nil, err
	}

	go func() {
		symbols, _, err := s.symbolService.GetPaged(ctx, &symbol_service.ReadPagedSymbolRequest{
			Filter: &symbol_service.SymbolFilter{
				PageSize:   20000,
				PageNumber: 1,
				Order:      "identifier",
				Ascending:  true,
			},
		})
		if err != nil {
			fmt.Printf("failed to get stored symbols %v", err)
			return
		}
		for _, sym := range *symbols {
			if err := stream.Send(sym); err != nil {
				fmt.Printf("failed while sending symbols to service: %v", err)
			}
		}
		if err := stream.CloseSend(); err != nil {
			fmt.Printf("failed to close send stream: %v", err)
			return
		}
	}()

	var result []*worker_symbol_service.RecalculateSymbolsResponse
	for {
		if res, err := stream.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error while receiving: %v", err)
			return nil, err
		} else {
			result = append(result, res)
		}
	}

	res, err := s.symbolService.ProcessRecalculationResponse(result, &ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
