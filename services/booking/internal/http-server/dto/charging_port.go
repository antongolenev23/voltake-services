package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type PortResponse struct {
	ID            uuid.UUID `json:"id"`
	ConnectorType string    `json:"connector_type"`
	PowerKW       int       `json:"power_kw"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewPortResponse(port domain.ChargingPort) PortResponse {
	return PortResponse{
		ID:            port.ID,
		ConnectorType: port.ConnectorType,
		PowerKW:       port.PowerKW,
		IsActive:      port.IsActive,
		CreatedAt:     port.CreatedAt,
	}
}

func NewPortsResponse(ports []domain.ChargingPort) []PortResponse {
	resp := make([]PortResponse, 0, len(ports))

	for _, port := range ports {
		resp = append(resp, NewPortResponse(port))
	}

	return resp
}

type PortRequest struct {
	ConnectorType string `json:"connector_type"`
	PowerKW       int    `json:"power_kw"`
	IsActive      bool   `json:"is_active"`
}

func (r PortRequest) ToCreateDomain(
	stationID uuid.UUID,
) domain.ChargingPort {
	return domain.ChargingPort{
		StationID:     stationID,
		ConnectorType: r.ConnectorType,
		PowerKW:       r.PowerKW,
		IsActive:      r.IsActive,
	}
}

func (r PortRequest) ToUpdateDomain(
	stationID uuid.UUID,
	portID uuid.UUID,
) domain.ChargingPort {
	return domain.ChargingPort{
		ID:            portID,
		StationID:     stationID,
		ConnectorType: r.ConnectorType,
		PowerKW:       r.PowerKW,
		IsActive:      r.IsActive,
	}
}
