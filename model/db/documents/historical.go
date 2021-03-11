package documents

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/vectorman1/analysis/analysis-api/generated/historical_service"
)

type Historical struct {
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

func (h *Historical) ToProtoObject() *historical_service.Historical {
	timestamp, _ := ptypes.TimestampProto(h.Timestamp)
	createdAt, _ := ptypes.TimestampProto(h.CreatedAt)

	return &historical_service.Historical{
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
