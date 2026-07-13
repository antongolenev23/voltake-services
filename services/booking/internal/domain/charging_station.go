package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChargingStation struct {
	ID        uuid.UUID
	Name      string
	Address   string
	Latitude  *float64
	Longitude *float64
	CreatedAt time.Time
}

func NewChargingStation(id uuid.UUID, name, address string) *ChargingStation {
	return &ChargingStation{
		ID:        id,
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
	}
}

func (s *ChargingStation) Rename(name string) {
	s.Name = name
}
