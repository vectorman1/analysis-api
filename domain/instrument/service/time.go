package service

import "time"

type TimeServiceContract interface {
	GetWorkdayTimeRange(start time.Time, workdays int) (from time.Time, to time.Time)
}

type TimeService struct {
}

func (t *TimeService) GetWorkdayTimeRange(start time.Time, workdays int) (from time.Time, to time.Time) {
	from = start
	to = start

	for i := 0; i < workdays; {
		switch from.Weekday() {
		case time.Sunday:
			from = from.Add(-1 * 24 * time.Hour)
		case time.Saturday:
			from = from.Add(-1 * 24 * time.Hour)
		default:
			from = from.Add(-1 * 24 * time.Hour)
			i++
		}
	}

	return from, to
}
