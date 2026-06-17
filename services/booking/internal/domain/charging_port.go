package domain

import "time"

type ChargingPort struct {
	ID            string
	StationID     string
	ConnectorType string
	PowerKW       int
	IsActive      bool
	CreatedAt     time.Time
}

func NewChargingPort(id, stationID, connectorType string, powerKW int) *ChargingPort {
	return &ChargingPort{
		ID:            id,
		StationID:     stationID,
		ConnectorType: connectorType,
		PowerKW:       powerKW,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}
}

func (p *ChargingPort) Activate() {
	p.IsActive = true
}

func (p *ChargingPort) Deactivate() {
	p.IsActive = false
}

func (p *ChargingPort) CanAcceptBooking() bool {
	return p.IsActive
}
