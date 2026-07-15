package domain

import "errors"

var (
	ErrInvalidTimeRange        = errors.New("invalid time range")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrBookingAlreadyCompleted = errors.New("booking already completed")
	ErrPortInactive            = errors.New("port is inactive")
	ErrStationNotFound         = errors.New("station not found")
	ErrPortNotFound            = errors.New("port not found")
)
