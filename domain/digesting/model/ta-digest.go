package model

import "time"

type TriggerType int

const (
	Lt TriggerType = iota
	Gt
	Rng
)

type TADigestRequest struct {
	InstrumentUuid  string
	ConsecutiveDays int
	TriggerType     TriggerType
	SourceProperty  string
	TargetProperty  string
	TargetNumber    int // used in case the trigger is a range from the source prop
	StartDate       time.Time
	EndDate         time.Time
}
