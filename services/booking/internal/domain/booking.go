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

	StartTime time.Time
	EndTime   time.Time

	Status BookingStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b Booking) Validate() error {
	if !b.EndTime.After(b.StartTime) {
		return ErrInvalidBookingPeriod
	}

	if b.StartTime.Before(time.Now()) {
		return ErrBookingInPast
	}

	duration := b.EndTime.Sub(b.StartTime)

	if duration < 30*time.Minute {
		return ErrBookingTooShort
	}

	if duration > 4*time.Hour {
		return ErrBookingTooLong
	}

	if b.StartTime.Minute()%30 != 0 ||
		b.EndTime.Minute()%30 != 0 {
		return ErrInvalidBookingTime
	}

	return nil
}

type BookingDetails struct {
	Booking

	Port    ChargingPort
	Station ChargingStation
}
