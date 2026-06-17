package domain

import (
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func NewTimeRange(start, end time.Time) (TimeRange, error) {
	if !end.After(start) {
		return TimeRange{}, ErrInvalidTimeRange
	}

	return TimeRange{
		Start: start,
		End:   end,
	}, nil
}
