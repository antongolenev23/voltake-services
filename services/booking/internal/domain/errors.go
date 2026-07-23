package domain

import "errors"

// Booking errors
var (
	ErrBookingAlreadyCancelled  = errors.New("booking already cancelled")
	ErrBookingAlreadyCompleted  = errors.New("booking already completed")
	ErrBookingConflict          = errors.New("booking conflict")
	ErrBookingInPast            = errors.New("booking in past")
	ErrBookingTooLong           = errors.New("booking too long")
	ErrBookingNotFound          = errors.New("booking not found")
	ErrBookingTooShort          = errors.New("booking duration is too short")
	ErrBookingCannotBeCancelled = errors.New("booking can not be cancelled")
	ErrInvalidBookingPeriod     = errors.New("booking end time must be after start time")
	ErrInvalidBookingTime       = errors.New("booking time must be aligned to 30 minutes")
	ErrInvalidTimeRange         = errors.New("invalid time range")
)

// Station errors
var (
	ErrStationNotFound = errors.New("station not found")
)

// Ports errors
var (
	ErrPortInactive    = errors.New("port is inactive")
	ErrPortNotFound    = errors.New("port not found")
	ErrPortUnavailable = errors.New("port unavailable")
)
