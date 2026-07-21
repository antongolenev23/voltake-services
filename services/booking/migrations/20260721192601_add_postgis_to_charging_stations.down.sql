DROP INDEX IF EXISTS idx_charging_stations_location;

ALTER TABLE charging_stations
ADD COLUMN latitude DOUBLE PRECISION,
ADD COLUMN longitude DOUBLE PRECISION;

ALTER TABLE charging_stations
DROP COLUMN location;

DROP EXTENSION IF EXISTS postgis;