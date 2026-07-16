package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChargingPort struct {
	ID            uuid.UUID
	StationID     uuid.UUID
	ConnectorType string
	PowerKW       int
	IsActive      bool
	CreatedAt     time.Time
}

func NewChargingPort(id, stationID uuid.UUID, connectorType string, powerKW int) *ChargingPort {
	return &ChargingPort{
		ID:            id,
		StationID:     stationID,
		ConnectorType: connectorType,
		PowerKW:       powerKW,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}
}

func (p *ChargingPort) CanAcceptBooking() bool {
	return p.IsActive
}
