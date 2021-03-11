package documents

import (
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/history_service"

	"github.com/golang/protobuf/ptypes"
)

type History struct {
	SymbolUuid string
	Open       float32
	Close      float32
	High       float32
	Low        float32
	Volume     int64
	AdjClose   float32
	Timestamp  time.Time
	CreatedAt  time.Time
}

func (h *History) ToProto() *history_service.History {
	timestamp, _ := ptypes.TimestampProto(h.Timestamp)
	createdAt, _ := ptypes.TimestampProto(h.CreatedAt)

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

func (a History) DistinctByTimestamp(in History) {

}
