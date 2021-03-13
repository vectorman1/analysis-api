package service

import (
	"context"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/vectorman1/analysis/analysis-api/third_party/yahoo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vectorman1/analysis/analysis-api/db"
)

type historyService interface {
	GetSymbolHistory(ctx context.Context, req *history_service.GetBySymbolUuidRequest) (*history_service.GetBySymbolUuidResponse, error)
	UpdateSymbolHistory(ctx context.Context, symUuid string, identifier string) (int, error)
	GetChartBySymbolUuid(ctx context.Context, req *history_service.GetChartBySymbolUuidRequest) (*history_service.GetChartBySymbolUuidResponse, error)
	UpdateAll(ctx context.Context) error
}

type HistoryService struct {
	yahooService             *yahoo.YahooService
	historyRepository        *db.HistoryRepository
	symbolRepository         *db.SymbolRepository
	symbolOverviewRepository *db.SymbolOverviewRepository
	reportService            *ReportService
}

func NewHistoryService(
	yahooService *yahoo.YahooService,
	historicalRepository *db.HistoryRepository,
	symbolRepository *db.SymbolRepository,
	symbolOverviewRepository *db.SymbolOverviewRepository,
	reportService *ReportService) *HistoryService {
	return &HistoryService{
		yahooService:             yahooService,
		historyRepository:        historicalRepository,
		symbolRepository:         symbolRepository,
		symbolOverviewRepository: symbolOverviewRepository,
		reportService:            reportService,
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
	// handle initial update of symbol
	if err != nil {
		beginningOfTime := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)

		// get history from Yahoo
		histories, err := s.yahooService.GetIdentifierHistory(
			symUuid,
			identifier,
			beginningOfTime,
			time.Now())
		if err != nil {
			return 0, err
		}

		// set Technical Analysis values based on histories
		*histories = s.reportService.GetTAValues(*histories, len(*histories))
		res, err := s.historyRepository.InsertMany(ctx, histories)
		if err != nil {
			return 0, err
		}

		return res, nil
	} else {
		// handle already existing history data
		// fetch history if last is older than 24h
		end := time.Now().UTC()
		if lastHistory.Timestamp.Add(time.Hour*24).Unix() > end.Unix() {
			return 0, nil
		}

		// get symbol history from (last + 24h) until now
		candles, err := s.yahooService.GetIdentifierHistory(
			symUuid,
			identifier,
			lastHistory.Timestamp.Add(time.Hour*24),
			end)
		if err != nil {
			return 0, err
		}

		if len(*candles) > 0 {
			// get last 150 history entries to pass to calculation method
			// it does further checks for each type of indicator
			// - e.g. min. 120 for MA120
			previous, err := s.historyRepository.GetSymbolHistory(ctx, symUuid, lastHistory.CreatedAt.Add(-(150 * (24 * time.Hour))), time.Now())
			if err != nil {
				return 0, err
			}

			// pass the new and old history for the calculation
			// the method returns only the difference
			*previous = append(*previous, *candles...)
			*candles = s.reportService.GetTAValues(*previous, len(*candles))

			res, err := s.historyRepository.InsertMany(ctx, candles)
			if err != nil {
				return 0, err
			}

			return res, nil
		}
	}

	return 0, nil
}

func (s *HistoryService) GetChartBySymbolUuid(ctx context.Context, req *history_service.GetChartBySymbolUuidRequest) (*history_service.GetChartBySymbolUuidResponse, error) {
	if !req.StartDate.IsValid() || !req.EndDate.IsValid() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date")
	}

	histories, err := s.historyRepository.GetSymbolHistory(ctx, req.SymbolUuid, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		return nil, err
	}

	var res history_service.GetChartBySymbolUuidResponse
	for _, h := range *histories {
		var values []float64
		// close open low high
		values = append(values, h.Close)
		values = append(values, h.Open)
		values = append(values, h.Low)
		values = append(values, h.High)

		res.Value = append(res.Value, &history_service.ChartValue{
			Date:   h.Timestamp.Format("2006-01-02"),
			Values: values,
		})
	}

	return &res, nil
}

func (s *HistoryService) UpdateAll(ctx context.Context) error {
	res, _, err := s.symbolRepository.GetPaged(
		context.Background(),
		&symbol_service.ReadPagedRequest{
			Filter: &symbol_service.SymbolFilter{
				PageSize:   100000,
				PageNumber: 1,
				Order:      "identifier",
				Ascending:  true,
			},
		})
	if err != nil {
		return err
	}

	grpclog.Infoln("[HISTORY JOB] Length of symbols to update: ", len(*res))

	hoursApprox := float32(len(*res) / 2000)
	grpclog.Infof("[HISTORY JOB] Job will take at least: %.2f hours", hoursApprox)

	for i, sym := range *res {
		// only update symbol history for these markets
		if sym.MarketName == "NASDAQ" ||
			sym.MarketName == "NYSE" {

			var u string
			sym.Uuid.AssignTo(&u)
			ctx, c := context.WithTimeout(ctx, 5*time.Second)
			entries, err := s.UpdateSymbolHistory(ctx, u, sym.Identifier)
			c()
			if err != nil {
				grpclog.Errorf("[HISTORY JOB] (%d/%d) Failed to update histories at: %s %s %s %s err: %v",
					i+1, len(*res),
					sym.Isin, sym.Identifier, sym.Name, sym.MarketName, err)
				continue
			} else if entries == 0 {
				grpclog.Infof("[HISTORY JOB] (%d/%d) No need to update: %s %s %s %s ",
					i+1, len(*res),
					sym.Isin, sym.Identifier, sym.Name, sym.MarketName)
				continue
			}

			grpclog.Infof("[HISTORY JOB] (%d/%d) Updated: %s %s %s %s Added entries: %d",
				i+1, len(*res),
				sym.Isin, sym.Identifier, sym.Name, sym.MarketName, entries)

			// timeout to avoid throttle
			time.Sleep(2 * time.Second)
		} else {
			grpclog.Infof("[HISTORY JOB] (%d/%d) Skipping: %s %s %s %s",
				i+1, len(*res),
				sym.Isin, sym.Identifier, sym.Name, sym.MarketName)
		}
	}

	return nil
}
