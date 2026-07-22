package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "booked"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID uuid.UUID

	UserID uuid.UUID
	PortID uuid.UUID

	StartTime     time.Time
	EndTime       time.Time
	ReservedUntil time.Time

	Status BookingStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b Booking) Validate(minDuration, maxDuration time.Duration) error {
	if !b.EndTime.After(b.StartTime) {
		return ErrInvalidBookingPeriod
	}

	if b.StartTime.Before(time.Now()) {
		return ErrBookingInPast
	}

	if b.StartTime.Second() != 0 ||
		b.StartTime.Nanosecond() != 0 ||
		b.EndTime.Second() != 0 ||
		b.EndTime.Nanosecond() != 0 {
		return ErrInvalidBookingTime
	}

	duration := b.EndTime.Sub(b.StartTime)

	if duration < minDuration {
		return ErrBookingTooShort
	}

	if duration > maxDuration {
		return ErrBookingTooLong
	}

	return nil
}

type BookingDetails struct {
	Booking

	Port    ChargingPort
	Station ChargingStation
}
