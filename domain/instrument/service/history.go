package service

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/repo"
	"github.com/vectorman1/analysis/analysis-api/domain/instrument/third_party"

	"github.com/vectorman1/analysis/analysis-api/common"

	validationErrors "github.com/vectorman1/analysis/analysis-api/common/errors"

	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HistoryServiceContract interface {
	GetSymbolHistory(ctx context.Context, req *instrument_service.HistoryRequest) (*instrument_service.HistoryResponse, error)
	UpdateSymbolHistory(ctx context.Context, symUuid string, identifier string) (int, error)
	GetChartBySymbolUuid(ctx context.Context, req *instrument_service.ChartRequest) (*instrument_service.ChartResponse, error)
	UpdateAll(ctx context.Context) error
}

type HistoryService struct {
	yahooService             *third_party.YahooService
	historyRepository        *repo.HistoryRepository
	symbolRepository         *repo.SymbolRepository
	symbolOverviewRepository *repo.SymbolOverviewRepository
	reportService            *ReportService
}

func NewHistoryService(
	yahooService *third_party.YahooService,
	historicalRepository *repo.HistoryRepository,
	symbolRepository *repo.SymbolRepository,
	symbolOverviewRepository *repo.SymbolOverviewRepository,
	reportService *ReportService) *HistoryService {
	return &HistoryService{
		yahooService:             yahooService,
		historyRepository:        historicalRepository,
		symbolRepository:         symbolRepository,
		symbolOverviewRepository: symbolOverviewRepository,
		reportService:            reportService,
	}
}

func (s *HistoryService) GetSymbolHistory(ctx context.Context, req *instrument_service.HistoryRequest) (*instrument_service.HistoryResponse, error) {
	if !req.StartDate.IsValid() || !req.EndDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid date range")
	}

	start := req.StartDate.AsTime()
	end := req.EndDate.AsTime()
	result, err := s.historyRepository.GetSymbolHistory(ctx, req.Uuid, start, end, true)
	if err != nil {
		return nil, err
	}

	var response []*instrument_service.History
	for _, history := range result {
		response = append(response, history.ToProto())
	}
	return &instrument_service.HistoryResponse{Items: response}, nil
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
		if len(*histories) == 0 {
			return 0, status.Error(codes.NotFound, validationErrors.NoHistoryFoundForSymbol)
		}

		// set Technical Analysis values based on histories
		*histories, _ = s.reportService.GetTAValues(*histories, len(*histories))
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
			previous, err := s.historyRepository.GetSymbolHistory(
				ctx,
				symUuid,
				lastHistory.Timestamp.Add(-(150 * (24 * time.Hour))),
				time.Now(),
				false)
			if err != nil {
				return 0, err
			}

			// pass the new and old history for the calculation
			// the method returns only the difference
			previous = append(previous, *candles...)
			*candles, err = s.reportService.GetTAValues(previous, len(*candles))
			if err != nil {
				return 0, err
			}

			res, err := s.historyRepository.InsertMany(ctx, candles)
			if err != nil {
				return 0, err
			}

			return res, nil
		}
	}

	return 0, nil
}

func (s *HistoryService) GetChartBySymbolUuid(
	ctx context.Context,
	req *instrument_service.ChartRequest) (*instrument_service.ChartResponse, error) {
	if !req.StartDate.IsValid() || !req.EndDate.IsValid() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date")
	}

	histories, err := s.historyRepository.GetSymbolHistory(
		ctx,
		req.SymbolUuid,
		req.StartDate.AsTime(),
		req.EndDate.AsTime(), false)
	if err != nil {
		return nil, err
	}

	if len(histories) == 0 || histories[len(histories)-1:][0].ShouldUpdate() {
		sym, err := s.symbolRepository.GetByUuid(ctx, req.SymbolUuid)
		if err != nil {
			return nil, status.Error(codes.NotFound, validationErrors.NoSymbolFound)
		}

		entries, err := s.UpdateSymbolHistory(ctx, req.SymbolUuid, sym.Identifier)
		if err != nil {
			return nil, err
		} else if entries == 0 {
			return nil, status.Error(codes.NotFound, validationErrors.NoHistoryFoundForSymbol)
		}

		if entries > 0 {
			histories, err = s.historyRepository.GetSymbolHistory(ctx, req.SymbolUuid, req.StartDate.AsTime(), req.EndDate.AsTime(), false)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, status.Error(codes.NotFound, validationErrors.NoHistoryFoundForSymbol)
		}
	}

	var res history_service.GetChartBySymbolUuidResponse
	for _, h := range histories {
		var values []float64
		// close open low high
		values = append(values, h.Open)
		values = append(values, h.Close)
		values = append(values, h.Low)
		values = append(values, h.High)

		res.Dates = append(res.Dates, h.Timestamp.Format("2006-01-02"))
		res.ChartDays = append(res.ChartDays, &history_service.ChartDay{Values: values})
	}

	return &res, nil
}

func (s *HistoryService) UpdateAll(ctx context.Context) error {
	res, _, err := s.symbolRepository.GetPaged(
		context.Background(),
		&symbol_service.GetPagedRequest{
			Filter: &proto_models.PagedFilter{
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

	processAvg := common.RollingAverage(10)

	for i, sym := range *res {
		start := time.Now()

		// only update symbol history for these markets
		if sym.MarketName == "NASDAQ" ||
			sym.MarketName == "NYSE" ||
			sym.MarketName == "OTC Markets" ||
			sym.MarketName == "NON-ISA OTC Markets" ||
			sym.MarketName == "NON-ISA NYSE" ||
			sym.MarketName == "NON-ISA NASDAQ" {

			var u string
			sym.Uuid.AssignTo(&u)
			ctx, c := context.WithTimeout(ctx, 10*time.Second)
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

		processTime := time.Since(start)
		newAvg := processAvg(processTime.Seconds())

		if i%25 == 0 {
			grpclog.Infof("[HISTORY JOB] MA of last 10 processed histories (per item): %2f seconds", newAvg)
			itemsLeft := len(*res) - i
			approxSecondsLeft := newAvg * float64(itemsLeft)
			grpclog.Infof("[HISTORY JOB] Estimated time left: %2f hours, or %2f minutes, or %2f seconds",
				approxSecondsLeft/60/60, approxSecondsLeft/60, approxSecondsLeft)
		}
	}

	return nil
}
