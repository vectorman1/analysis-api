package service

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
)

type CandleType uint8

const (
	Green CandleType = iota
	Red
)

type daily interface {
	LastNDays(candleType CandleType, n int) bool
}

type Daily struct {
}

func (d *Daily) LastNDays(candleType CandleType, n int, iter *chart.Iter) (bool, error) {
	for i := 0; i < n; i++ {
		if iter.Next() {
			bar := iter.Bar()
			barCandleType := d.getCandleType(bar)

			if barCandleType != candleType {
				return false, nil
			}
		} else {
			return false, fmt.Errorf("error while getting price history")
		}
	}

	return true, nil
}

func (d *Daily) getCandleType(bar *finance.ChartBar) CandleType {
	if bar.AdjClose.GreaterThan(bar.Open) {
		return Green
	}

	return Red
}
