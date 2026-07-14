DROP INDEX IF EXISTS idx_charging_stations_owner_id;

ALTER TABLE charging_stations
DROP COLUMN owner_id;