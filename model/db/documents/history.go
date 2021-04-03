package documents

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"
)

type MA struct {
	MA5   float64
	MA10  float64
	MA20  float64
	MA30  float64
	MA60  float64
	MA120 float64
}

type EMA struct {
	EMA5   float64
	EMA10  float64
	EMA20  float64
	EMA30  float64
	EMA60  float64
	EMA120 float64
}

type MACD struct {
	Line      float64
	Histogram float64
}

type Trend struct {
	Trend5   float64
	Trend10  float64
	Trend20  float64
	Trend30  float64
	Trend60  float64
	Trend120 float64
}

type LastHistory struct {
	Close     float64
	Timestamp time.Time
}

type History struct {
	SymbolUuid string
	Calculated bool
	Open       float64
	Close      float64
	High       float64
	Low        float64
	Volume     int64
	AdjClose   float64

	Trend Trend
	MA    MA
	EMA   EMA
	MACD  MACD
	RSI   float64

	Timestamp time.Time
	CreatedAt time.Time
}

func (h *History) ToProto() *history_service.History {
	timestamp := timestamppb.New(h.Timestamp)
	createdAt := timestamppb.New(h.CreatedAt)

	return &history_service.History{
		Open:      h.Open,
		Close:     h.Close,
		High:      h.High,
		Low:       h.Low,
		Volume:    h.Volume,
		AdjClose:  h.AdjClose,
		Timestamp: timestamp,
		CreatedAt: createdAt,
	}
}
