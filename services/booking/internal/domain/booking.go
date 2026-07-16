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

type BookingDetails struct {
	Booking

	Port    ChargingPort
	Station ChargingStation
}
