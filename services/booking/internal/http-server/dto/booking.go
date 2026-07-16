package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type BookingRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (r BookingRequest) ToDomain(userID, portID uuid.UUID) domain.Booking {
	return domain.Booking{
		UserID:    userID,
		PortID:    portID,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	}
}

type BookingResponse struct {
	ID        uuid.UUID            `json:"id"`
	UserID    uuid.UUID            `json:"user_id"`
	PortID    uuid.UUID            `json:"port_id"`
	StartTime time.Time            `json:"start_time"`
	EndTime   time.Time            `json:"end_time"`
	Status    domain.BookingStatus `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

func NewBookingResponse(b domain.Booking) BookingResponse {
	return BookingResponse{
		ID:        b.ID,
		UserID:    b.UserID,
		PortID:    b.PortID,
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

type BookingListResponse struct {
	ID        uuid.UUID            `json:"id"`
	PortID    uuid.UUID            `json:"port_id"`
	StartTime time.Time            `json:"start_time"`
	EndTime   time.Time            `json:"end_time"`
	Status    domain.BookingStatus `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
}

func NewBookingListResponse(
	bookings []domain.Booking,
) []BookingListResponse {
	resp := make([]BookingListResponse, 0, len(bookings))

	for _, b := range bookings {
		resp = append(resp, BookingListResponse{
			ID:        b.ID,
			PortID:    b.PortID,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			Status:    b.Status,
			CreatedAt: b.CreatedAt,
		})
	}

	return resp
}

type BookingDetailsResponse struct {
	ID        uuid.UUID            `json:"id"`
	StartTime time.Time            `json:"start_time"`
	EndTime   time.Time            `json:"end_time"`
	Status    domain.BookingStatus `json:"status"`

	Port PortDetailsResponse `json:"port"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PortDetailsResponse struct {
	PortResponse

	Station StationResponse `json:"station"`
}

func NewBookingDetailsResponse(
	b domain.BookingDetails,
) BookingDetailsResponse {
	return BookingDetailsResponse{
		ID:        b.ID,
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,

		Port: PortDetailsResponse{
			PortResponse: PortResponse{
				ID:            b.Port.ID,
				ConnectorType: b.Port.ConnectorType,
				PowerKW:       b.Port.PowerKW,
				IsActive:      b.Port.IsActive,
				CreatedAt:     b.Port.CreatedAt,
			},

			Station: StationResponse{
				ID:        b.Station.ID,
				Name:      b.Station.Name,
				Address:   b.Station.Address,
				Latitude:  b.Station.Latitude,
				Longitude: b.Station.Longitude,
			},
		},
	}
}
