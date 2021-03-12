package documents

import "time"

type valueType uint

const (
	Indicator valueType = iota
	Constant
)

type operator uint

const (
	LessThan operator = iota
	MoreThan
	MoreThanOrEqualTo
	LessThanOrEqualTo
)

type Trigger struct {
	UserID          string
	IsConsecutive   bool
	Consecutive     uint
	SourceValueType valueType
	TargetValueType valueType

	SourceValueKey string
	TargetValueKey string

	Operator operator

	UpdatedAt time.Time
	CreatedAt time.Time
}
