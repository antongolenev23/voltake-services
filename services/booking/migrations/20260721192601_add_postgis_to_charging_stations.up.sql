CREATE EXTENSION IF NOT EXISTS postgis;

ALTER TABLE charging_stations
ADD COLUMN location geography(Point, 4326);

ALTER TABLE charging_stations
DROP COLUMN latitude,
DROP COLUMN longitude;

CREATE INDEX idx_charging_stations_location
ON charging_stations
USING GIST(location);