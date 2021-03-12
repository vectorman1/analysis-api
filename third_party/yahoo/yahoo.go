package yahoo

import (
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/vectorman1/analysis/analysis-api/model/db/documents"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type yahooService interface {
	GetIdentifierHistory(symUuid string, identifier string, start time.Time, end time.Time) (*[]documents.History, error)
}

type YahooService struct {
	yahooService
}

func NewYahooService() *YahooService {
	return &YahooService{}
}

func (s *YahooService) GetIdentifierHistory(symUuid string, identifier string, start time.Time, end time.Time) (*[]documents.History, error) {
	params := &chart.Params{
		Symbol:   identifier,
		Interval: datetime.OneDay,
		End:      datetime.FromUnix(int(end.Unix())),
		Start:    datetime.FromUnix(int(start.Unix())),
	}
	iter := chart.Get(params)

	var result []documents.History
	for iter.Next() {
		bar := iter.Bar()

		open, _ := bar.Open.Float64()
		cl, _ := bar.Close.Float64()
		high, _ := bar.High.Float64()
		low, _ := bar.Low.Float64()
		adjClose, _ := bar.AdjClose.Float64()
		timestamp := time.Unix(int64(bar.Timestamp), 0)

		result = append(result, documents.History{
			SymbolUuid: symUuid,
			Open:       open,
			Close:      cl,
			High:       high,
			Low:        low,
			Volume:     int64(bar.Volume),
			AdjClose:   adjClose,
			Timestamp:  timestamp,
			CreatedAt:  time.Now(),
		})
	}

	if err := iter.Err(); err != nil {
		return nil, status.Error(codes.FailedPrecondition, "new history data was null")
	}

	return &result, nil
}
