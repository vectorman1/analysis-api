package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/vectorman1/analysis/analysis-api/db"
)

type historicalService interface {
}

type HistoricalService struct {
	historicalService
	historicalRepository *db.HistoricalRepository
	symbolRepository     *db.SymbolRepository
}

func NewHistoricalService(historicalRepository *db.HistoricalRepository, symbolRepository *db.SymbolRepository) *HistoricalService {
	return &HistoricalService{
		historicalRepository: historicalRepository,
		symbolRepository:     symbolRepository,
	}
}

func (h *HistoricalService) GetBySymbolUuid(ctx *context.Context, req *historical_service.GetBySymbolUuidRequest) (*[]*historical_service.Historical, error) {
	sym, err := h.symbolRepository.GetByUuid(*ctx, req.SymbolUuid)
	if err != nil {
		return nil, err
	}

	marketAbbr := ""
	switch marketAbbr {
	case "Bolsa de Madrid":
		marketAbbr = "BME"
		break
	case "Deutsche Börse Xetra":
		marketAbbr = "XETR"
		break
	case "Euronext Paris":
		marketAbbr = "XPAR"
		break
	case "London Stock Exchange":
		marketAbbr = "XLON"
		break
	case "LSE AIM":
		marketAbbr = "AIMX"
		break
	}

	params := &chart.Params{
		Symbol:   sym.Identifier,
		Interval: datetime.OneDay,
		End:      datetime.FromUnix(int(time.Now().Unix())),
		Start:    datetime.FromUnix(int(time.Date(1970, time.Month(1), 1, 0, 0, 0, 0, time.UTC).Unix())),
	}
	iter := chart.Get(params)

	var result []*historical_service.Historical
	for iter.Next() {
		bar := iter.Bar()
		date := time.Unix(int64(bar.Timestamp), 0)
		println(date.String())

		result = append(result)
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}

	return nil, nil
}
