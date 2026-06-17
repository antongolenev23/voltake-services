package domain

import "time"

type ChargingStation struct {
	ID        string
	Name      string
	Address   string
	Latitude  *float64
	Longitude *float64
	CreatedAt time.Time
}

func NewChargingStation(id, name, address string) *ChargingStation {
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
