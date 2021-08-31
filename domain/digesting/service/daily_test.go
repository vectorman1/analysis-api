package service

import (
	"testing"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type MockSuite struct {
	Daily
}

func TestLastNDaysOfSpecificCandle(t *testing.T) {
	endDate := time.Date(2021, 7, 9, 23, 0, 0, 0, time.UTC)
	startDate := endDate.Add(-10 * 24 * time.Hour) // needs to be 2 extra days to account for weekend

	params := &chart.Params{
		Symbol:   "AAPL",
		Interval: datetime.OneDay,
		Start:    datetime.New(&startDate),
		End:      datetime.New(&endDate),
	}

	iter := chart.Get(params)

	mockSuite := &MockSuite{Daily{}}

	ok, err := mockSuite.LastNDays(Green, 8, iter)

	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
