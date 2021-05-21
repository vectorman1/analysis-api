package present

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"

	service2 "github.com/vectorman1/analysis/analysis-api/domain/instrument/service"

	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/jobs"

	"github.com/vectorman1/analysis/analysis-api/common"
)

type InstrumentServiceServer struct {
	rabbitClient   *common.RabbitClient
	symbolService  *service2.InstrumentsService
	historyService *service2.HistoryService
	instrument_service.UnimplementedInstrumentServiceServer
}

func NewSymbolServiceServer(symbolsService *service2.InstrumentsService, historyService *service2.HistoryService) *InstrumentServiceServer {
	return &InstrumentServiceServer{
		symbolService:  symbolsService,
		historyService: historyService,
	}
}

func (s *InstrumentServiceServer) Get(
	ctx context.Context,
	req *instrument_service.InstrumentRequest) (*instrument_service.Instrument, error) {
	return s.symbolService.Get(ctx, req.Uuid)
}

func (s *InstrumentServiceServer) GetPaged(
	ctx context.Context,
	req *instrument_service.PagedRequest) (*instrument_service.PagedResponse, error) {
	timeoutContext, c := context.WithTimeout(ctx, 5*time.Second)
	defer c()

	res, totalItemsCount, err := s.symbolService.GetPaged(timeoutContext, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	resp := &instrument_service.PagedResponse{
		Items:      *res,
		TotalItems: uint64(totalItemsCount),
	}
	return resp, nil
}

func (s *InstrumentServiceServer) Overview(
	ctx context.Context,
	req *instrument_service.InstrumentRequest) (*instrument_service.InstrumentOverview, error) {
	res, err := s.symbolService.Overview(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return res, nil
}

func (s *InstrumentServiceServer) UpdateAllJob(
	ctx context.Context,
	req *instrument_service.StartUpdateJobRequest) (*instrument_service.StartUpdateJobResponse, error) {
	jobrunner.Now(jobs.NewSymbolUpdateJob(s.symbolService))

	return &instrument_service.StartUpdateJobResponse{}, nil
}

func (s *InstrumentServiceServer) UpdateAll(
	ctx context.Context,
	req *instrument_service.StartUpdateJobRequest) (*instrument_service.UpdateAllResponse, error) {
	res, err := s.symbolService.UpdateAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
