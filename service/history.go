package service

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/vectorman1/analysis/analysis-api/third_party/yahoo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vectorman1/analysis/analysis-api/db"
)

type historyService interface {
	GetSymbolHistory(ctx context.Context, req *history_service.GetBySymbolUuidRequest) (*history_service.GetBySymbolUuidResponse, error)
	UpdateSymbolHistory(ctx context.Context, symUuid string, identifier string) (int, error)
}

type HistoryService struct {
	historyService
	yahooService             *yahoo.YahooService
	historyRepository        *db.HistoryRepository
	symbolRepository         *db.SymbolRepository
	symbolOverviewRepository *db.SymbolOverviewRepository
}

func NewHistoryService(
	yahooService *yahoo.YahooService,
	historicalRepository *db.HistoryRepository,
	symbolRepository *db.SymbolRepository,
	symbolOverviewRepository *db.SymbolOverviewRepository) *HistoryService {
	return &HistoryService{
		yahooService:             yahooService,
		historyRepository:        historicalRepository,
		symbolRepository:         symbolRepository,
		symbolOverviewRepository: symbolOverviewRepository,
	}
}

func (s *HistoryService) GetSymbolHistory(ctx context.Context, req *history_service.GetBySymbolUuidRequest) (*history_service.GetBySymbolUuidResponse, error) {
	if !req.StartDate.IsValid() || !req.EndDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid date range")
	}

	start := req.StartDate.AsTime()
	end := req.EndDate.AsTime()
	result, err := s.historyRepository.GetSymbolHistory(ctx, req.SymbolUuid, start, end)
	if err != nil {
		return nil, err
	}

	var response []*history_service.History
	for _, history := range *result {
		response = append(response, history.ToProto())
	}
	return &history_service.GetBySymbolUuidResponse{Items: response}, nil
}

func (s *HistoryService) UpdateSymbolHistory(ctx context.Context, symUuid string, identifier string) (int, error) {
	lastHistory, err := s.historyRepository.GetLastSymbolHistory(ctx, symUuid)
	if err != nil {
		return 0, err
	}

	// handle initial update of symbol
	if lastHistory == nil {
		beginningOfTime := time.Unix(0, 0)

		candles, err := s.yahooService.GetIdentifierHistory(
			symUuid,
			identifier,
			beginningOfTime,
			time.Now())
		if err != nil {
			return 0, err
		}

		res, err := s.historyRepository.InsertMany(ctx, *candles)
		if err != nil {
			return 0, err
		}

		return res, nil
	} else {
		// handle already existing history data
		// fetch fetch history if last is older than 24h
		end := time.Now().UTC()
		if lastHistory.Timestamp.Add(time.Hour*24).Unix() > end.Unix() {
			return 0, nil
		}

		candles, err := s.yahooService.GetIdentifierHistory(
			symUuid,
			identifier,
			lastHistory.Timestamp.Add(time.Hour*24),
			end)
		if err != nil {
			return 0, err
		}

		if len(*candles) > 0 {
			res, err := s.historyRepository.InsertMany(ctx, *candles)
			if err != nil {
				return 0, err
			}

			return res, nil
		}
	}

	return 0, nil
}
