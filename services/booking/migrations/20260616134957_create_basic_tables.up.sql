CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE booking_status AS ENUM (
    'booked',
    'cancelled',
    'completed'
);

CREATE TABLE charging_stations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE charging_ports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id UUID NOT NULL REFERENCES charging_stations(id) ON DELETE CASCADE,
    connector_type TEXT NOT NULL,
    power_kw INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_charging_ports_station_id
ON charging_ports(station_id);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    port_id UUID NOT NULL REFERENCES charging_ports(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status booking_status NOT NULL DEFAULT 'booked',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_user_id
ON bookings(user_id);

CREATE INDEX idx_bookings_port_id
ON bookings(port_id);

CREATE INDEX idx_bookings_port_time
ON bookings(port_id, start_time, end_time);