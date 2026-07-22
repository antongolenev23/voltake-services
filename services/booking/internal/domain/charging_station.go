package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrStationNameEmpty        = errors.New("station name is empty")
	ErrStationAddressEmpty     = errors.New("station address is empty")
	ErrStationInvalidLatitude  = errors.New("invalid station latitude")
	ErrStationInvalidLongitude = errors.New("invalid station longitude")
	ErrInvalidOwnerID          = errors.New("invalid owner id")
)

type ChargingStation struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
}

type ChargingStationDetails struct {
	ChargingStation
	Ports []ChargingPort
}

func NewChargingStation(id uuid.UUID, name, address string) *ChargingStation {
	return &ChargingStation{
		ID:        id,
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
	}
}

func (s *ChargingStation) Validate() error {
	if s.Name == "" {
		return ErrStationNameEmpty
	}

	if s.OwnerID == uuid.Nil {
		return ErrInvalidOwnerID
	}

	if s.Address == "" {
		return ErrStationAddressEmpty
	}

	if s.Latitude < -90 || s.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	if s.Longitude < -180 || s.Longitude > 180 {
		return errors.New("invalid longitude")
	}

	return nil
}
