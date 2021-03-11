package service

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/model/db/documents"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/vectorman1/analysis/analysis-api/db"
)

type historicalService interface {
	GetBySymbolUuid(ctx context.Context, req *historical_service.GetBySymbolUuidRequest) (*historical_service.GetBySymbolUuidResponse, error)
}

type HistoricalService struct {
	historicalService
	historicalRepository     *db.HistoricalRepository
	symbolRepository         *db.SymbolRepository
	symbolOverviewRepository *db.SymbolOverviewRepository
}

func NewHistoricalService(historicalRepository *db.HistoricalRepository, symbolRepository *db.SymbolRepository, symbolOverviewRepository *db.SymbolOverviewRepository) *HistoricalService {
	return &HistoricalService{
		historicalRepository:     historicalRepository,
		symbolRepository:         symbolRepository,
		symbolOverviewRepository: symbolOverviewRepository,
	}
}

func (h *HistoricalService) GetBySymbolUuid(ctx context.Context, req *historical_service.GetBySymbolUuidRequest) (*historical_service.GetBySymbolUuidResponse, error) {
	if !req.StartDate.IsValid() || !req.EndDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid date range")
	}

	start := req.StartDate.AsTime()
	end := req.EndDate.AsTime()
	result, err := h.historicalRepository.GetBySymbolUuid(ctx, req.SymbolUuid, start, end)
	if err != nil {
		return nil, err
	}

	if result == nil {
		sym, err := h.symbolRepository.GetByUuid(ctx, req.SymbolUuid)
		if err != nil {
			return nil, err
		}

		params := &chart.Params{
			Symbol:   sym.Identifier,
			Interval: datetime.OneDay,
			End:      datetime.FromUnix(int(time.Now().Unix())),
			Start:    datetime.FromUnix(int(time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC).Unix())),
		}
		iter := chart.Get(params)

		for iter.Next() {
			bar := iter.Bar()

			open, _ := bar.Open.Float64()
			cl, _ := bar.Close.Float64()
			high, _ := bar.High.Float64()
			low, _ := bar.Low.Float64()
			adjClose, _ := bar.AdjClose.Float64()
			timestamp := time.Unix(int64(bar.Timestamp), 0)

			*result = append(*result, documents.Historical{
				SymbolUuid: req.SymbolUuid,
				Open:       float32(open),
				Close:      float32(cl),
				High:       float32(high),
				Low:        float32(low),
				Volume:     int64(bar.Volume),
				AdjClose:   float32(adjClose),
				Timestamp:  timestamp,
				CreatedAt:  time.Now(),
			})
		}
		if err := iter.Err(); err != nil {
			return nil, status.Error(codes.FailedPrecondition, "new history data was null")
		}

		_, err = h.historicalRepository.InsertMany(ctx, *result)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed saving histories: %v", err)
		}
	}

	var response []*historical_service.Historical
	for _, history := range *result {
		response = append(response, history.ToProtoObject())
	}
	return &historical_service.GetBySymbolUuidResponse{Items: response}, nil
}
