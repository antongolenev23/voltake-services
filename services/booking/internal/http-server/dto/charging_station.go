package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type StationRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (r *StationRequest) ToCreateDomain(
	ownerID uuid.UUID,
) domain.ChargingStation {
	return domain.ChargingStation{
		OwnerID:   ownerID,
		Name:      r.Name,
		Address:   r.Address,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

func (r *StationRequest) ToUpdateDomain(
	id uuid.UUID,
	ownerID uuid.UUID,
) domain.ChargingStation {
	return domain.ChargingStation{
		ID:        id,
		OwnerID:   ownerID,
		Name:      r.Name,
		Address:   r.Address,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

type StationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewStationResponse(station domain.ChargingStation) StationResponse {
	return StationResponse{
		ID:        station.ID,
		Name:      station.Name,
		Address:   station.Address,
		Latitude:  station.Latitude,
		Longitude: station.Longitude,
		CreatedAt: station.CreatedAt,
	}
}

func NewStationsResponse(
	stations []domain.ChargingStation,
) []StationResponse {

	responses := make([]StationResponse, 0, len(stations))

	for _, station := range stations {
		responses = append(
			responses,
			NewStationResponse(station),
		)
	}

	return responses
}
