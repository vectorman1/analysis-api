package service

import (
	"sync"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sdcoffey/big"

	"github.com/sdcoffey/techan"
)

type reportService interface {
	GetTAValues(history *[]model.History) *[]model.History
}

type ReportService struct {
}

func NewReportService() *ReportService {
	return &ReportService{}
}

// GetTAValues calculates the respective TA values of the histories and returns
// a sub-slice with len newLen, starting from the end.
func (s *ReportService) GetTAValues(histories []model.History, newLen int) ([]model.History, error) {
	historiesLen := len(histories)
	if historiesLen < 5 {
		return nil, nil
	}

	series := techan.NewTimeSeries()
	for _, history := range histories {
		period := techan.NewTimePeriod(history.Timestamp, 24*time.Hour)
		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewDecimal(history.Open)
		candle.ClosePrice = big.NewDecimal(history.Close)
		candle.MaxPrice = big.NewDecimal(history.High)
		candle.MinPrice = big.NewDecimal(history.Low)
		candle.Volume = big.NewDecimal(float64(history.Volume))
		if ok := series.AddCandle(candle); !ok {
			return nil, status.Error(codes.FailedPrecondition, "History data error")
		}
	}

	closePrices := techan.NewClosePriceIndicator(series)
	ma5 := techan.NewSimpleMovingAverage(closePrices, 5)
	ma10 := techan.NewSimpleMovingAverage(closePrices, 10)
	ma20 := techan.NewSimpleMovingAverage(closePrices, 20)
	ma30 := techan.NewSimpleMovingAverage(closePrices, 30)
	ma60 := techan.NewSimpleMovingAverage(closePrices, 60)
	ma120 := techan.NewSimpleMovingAverage(closePrices, 120)
	ema5 := techan.NewEMAIndicator(closePrices, 5)
	ema10 := techan.NewEMAIndicator(closePrices, 10)
	ema20 := techan.NewEMAIndicator(closePrices, 20)
	ema30 := techan.NewEMAIndicator(closePrices, 30)
	ema60 := techan.NewEMAIndicator(closePrices, 60)
	ema120 := techan.NewEMAIndicator(closePrices, 120)
	trend5 := techan.NewTrendlineIndicator(closePrices, 5)
	trend10 := techan.NewTrendlineIndicator(closePrices, 10)
	trend20 := techan.NewTrendlineIndicator(closePrices, 20)
	trend30 := techan.NewTrendlineIndicator(closePrices, 30)
	trend60 := techan.NewTrendlineIndicator(closePrices, 60)
	trend120 := techan.NewTrendlineIndicator(closePrices, 120)
	macd := techan.NewMACDIndicator(closePrices, 12, 26)
	macdHist := techan.NewMACDHistogramIndicator(closePrices, 9)
	rsi := techan.NewRelativeStrengthIndexIndicator(closePrices, 9)

	var wg sync.WaitGroup
	for i, history := range histories {
		if history.Calculated {
			continue
		}

		wg.Add(1)
		go func(i int, ma5, ma10, ma20, ma30, ma60, ma120 techan.Indicator) {
			defer wg.Done()
			histories[i].MA.MA5 = ma5.Calculate(i).Float()
			histories[i].MA.MA10 = ma10.Calculate(i).Float()
			histories[i].MA.MA20 = ma20.Calculate(i).Float()
			histories[i].MA.MA30 = ma30.Calculate(i).Float()
			histories[i].MA.MA60 = ma60.Calculate(i).Float()
			histories[i].MA.MA120 = ma120.Calculate(i).Float()
		}(i, ma5, ma10, ma20, ma30, ma60, ma120)

		wg.Add(1)
		go func(i int, ema5, ema10, ema20, ema30, ema60, ema120 techan.Indicator) {
			defer wg.Done()
			histories[i].EMA.EMA5 = ema5.Calculate(i).Float()
			histories[i].EMA.EMA10 = ema10.Calculate(i).Float()
			histories[i].EMA.EMA20 = ema20.Calculate(i).Float()
			histories[i].EMA.EMA30 = ema30.Calculate(i).Float()
			histories[i].EMA.EMA60 = ema60.Calculate(i).Float()
			histories[i].EMA.EMA120 = ema120.Calculate(i).Float()
		}(i, ema5, ema10, ema20, ema30, ema60, ema120)

		wg.Add(1)
		go func(i int, trend5, trend10, trend20, trend30, trend60, trend120 techan.Indicator) {
			defer wg.Done()
			if i >= 5 {
				histories[i].Trend.Trend5 = trend5.Calculate(i).Float()
			}
			if i >= 10 {
				histories[i].Trend.Trend10 = trend10.Calculate(i).Float()
			}
			if i >= 20 {
				histories[i].Trend.Trend20 = trend20.Calculate(i).Float()
			}
			if i >= 30 {
				histories[i].Trend.Trend30 = trend30.Calculate(i).Float()
			}
			if i >= 60 {
				histories[i].Trend.Trend60 = trend60.Calculate(i).Float()
			}
			if i >= 120 {
				histories[i].Trend.Trend120 = trend120.Calculate(i).Float()
			}
		}(i, trend5, trend10, trend20, trend30, trend60, trend120)

		wg.Add(1)
		go func(i int, macd, macdHist, rsi techan.Indicator) {
			defer wg.Done()
			histories[i].MACD.Line = macd.Calculate(i).Float()
			histories[i].MACD.Histogram = macdHist.Calculate(i).Float()
			histories[i].RSI = rsi.Calculate(i).Float()
		}(i, macd, macdHist, rsi)

		histories[i].Calculated = true
	}

	wg.Wait()

	return histories[historiesLen-newLen:], nil
}
