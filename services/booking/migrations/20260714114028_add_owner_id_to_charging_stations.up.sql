ALTER TABLE charging_stations
ADD COLUMN owner_id UUID NOT NULL;

CREATE INDEX idx_charging_stations_owner_id
ON charging_stations(owner_id);