package model

import "time"

type TriggerType int

const (
	Lt TriggerType = iota
	Gt
	Rng
)

type TADigestRequest struct {
	SymbolUuid      string
	ConsecutiveDays int
	TriggerType     TriggerType
	SourceProperty  string
	TargetProperty  string
	TargetNumber    int
	StartDate       time.Time
	EndDate         time.Time
}
