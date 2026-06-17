package domain

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "booked"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID     string
	UserID string
	PortID string

	Range  TimeRange
	Status BookingStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBooking(id, userID, portID string, r TimeRange) *Booking {
	now := time.Now()

	return &Booking{
		ID:        id,
		UserID:    userID,
		PortID:    portID,
		Range:     r,
		Status:    BookingStatusBooked,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Booking) Cancel() {
	if b.Status != BookingStatusBooked {
		return
	}
	b.Status = BookingStatusCancelled
	b.UpdatedAt = time.Now()
}

func (b *Booking) Complete() {
	if b.Status != BookingStatusBooked {
		return
	}
	b.Status = BookingStatusCompleted
	b.UpdatedAt = time.Now()
}

func (b *Booking) IsActive() bool {
	return b.Status == BookingStatusBooked
}
