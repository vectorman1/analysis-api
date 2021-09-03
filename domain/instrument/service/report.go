package service

import (
	"sync"
	"time"

	"github.com/vectorman1/analysis/analysis-api/common/errors"

	"go.uber.org/zap"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"
	logger_grpc "github.com/vectorman1/analysis/analysis-api/middleware/logger-grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sdcoffey/big"

	"github.com/sdcoffey/techan"
)

type reportService interface {
	GetTAValues(history *[]model.History) *[]model.History
}

type TAService struct {
}

func NewReportService() *TAService {
	return &TAService{}
}

// GetTAValues calculates the respective TA values of the histories and returns
// a sub-slice with len newLen, starting from the end.
func (s *TAService) GetTAValues(histories []model.History, newLen int) ([]model.History, error) {
	historiesLen := len(histories)
	if historiesLen < 5 {
		logger_grpc.Log.Debug("histores length is less than 5 so no TA can be done",
			zap.Int("histories-len", historiesLen))
		return nil, nil
	}

	series := techan.NewTimeSeries()
	for _, history := range histories {
		period := techan.NewTimePeriod(history.Timestamp, 24*time.Hour)
		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewDecimal(history.Open)
		candle.ClosePrice = big.NewDecimal(history.AdjClose)
		candle.MaxPrice = big.NewDecimal(history.High)
		candle.MinPrice = big.NewDecimal(history.Low)
		candle.Volume = big.NewDecimal(float64(history.Volume))
		if ok := series.AddCandle(candle); !ok {
			logger_grpc.Log.Debug(errors.FailedToAddCandleToTimeSeries,
				zap.String("attempted-candle", candle.Period.String()),
				zap.String("last-time-series-entry", series.LastCandle().Period.String()))
			return nil, status.Error(codes.Internal, errors.FailedToAddCandleToTimeSeries)
		}
	}

	seriesLen := len(series.Candles)
	if seriesLen != historiesLen {
		logger_grpc.Log.Debug(errors.OutputtedTimeSeriesMismatchedInputHistoriesLen,
			zap.Int("input-histories-len", historiesLen),
			zap.Int("output-time-series-len", seriesLen))
		return nil, status.Error(codes.Internal, errors.OutputtedTimeSeriesMismatchedInputHistoriesLen)
	}

	closePrices := techan.NewClosePriceIndicator(series)

	// Create Technical analysis calculation structs
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

		// Different lengths of minimum time period is needed for different TA indicators
		if i >= 5 {
			wg.Add(1)
			go func(i int, ma5, ema5, trend5 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA5 = ma5.Calculate(i).Float()
				histories[i].EMA.EMA5 = ema5.Calculate(i).Float()
				histories[i].Trend.Trend5 = trend5.Calculate(i).Float()
			}(i, ma5, ema5, trend5)
		}

		if i >= 9 {
			wg.Add(1)
			go func(i int, macdHist techan.Indicator) {
				defer wg.Done()
				histories[i].MACD.Histogram = macdHist.Calculate(i).Float()
				histories[i].RSI = rsi.Calculate(i).Float()
			}(i, macdHist)
		}

		if i >= 10 {
			wg.Add(1)
			go func(i int, ma10, ema10, trend10 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA10 = ma10.Calculate(i).Float()
				histories[i].EMA.EMA10 = ema10.Calculate(i).Float()
				histories[i].Trend.Trend10 = trend10.Calculate(i).Float()
			}(i, ma10, ema10, trend10)
		}

		if i >= 20 {
			wg.Add(1)
			go func(i int, ma20, ema20, trend20 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA20 = ma20.Calculate(i).Float()
				histories[i].EMA.EMA20 = ema20.Calculate(i).Float()
				histories[i].Trend.Trend20 = trend20.Calculate(i).Float()
			}(i, ma20, ema20, trend20)
		}

		if i >= 26 {
			wg.Add(1)
			go func(i int, macd techan.Indicator) {
				defer wg.Done()
				histories[i].MACD.Line = macd.Calculate(i).Float()
			}(i, macd)
		}

		if i >= 30 {
			wg.Add(1)
			go func(i int, ma30, ema30, trend30 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA30 = ma30.Calculate(i).Float()
				histories[i].EMA.EMA30 = ema30.Calculate(i).Float()
				histories[i].Trend.Trend30 = trend30.Calculate(i).Float()
			}(i, ma30, ema30, trend30)
		}

		if i >= 60 {
			wg.Add(1)
			go func(i int, ma60, ema60, trend60 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA60 = ma60.Calculate(i).Float()
				histories[i].EMA.EMA60 = ema60.Calculate(i).Float()
				histories[i].Trend.Trend60 = trend60.Calculate(i).Float()
			}(i, ma60, ema60, trend60)
		}

		if i >= 120 {
			wg.Add(1)
			go func(i int, ma120, ema120, trend120 techan.Indicator) {
				defer wg.Done()
				histories[i].MA.MA120 = ma120.Calculate(i).Float()
				histories[i].EMA.EMA120 = ema120.Calculate(i).Float()
				histories[i].Trend.Trend120 = trend120.Calculate(i).Float()
			}(i, ma120, ema120, trend120)
		}

		histories[i].Calculated = true
	}

	wg.Wait()

	return histories[historiesLen-newLen:], nil
}
