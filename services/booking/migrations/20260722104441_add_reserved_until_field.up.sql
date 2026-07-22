CREATE EXTENSION IF NOT EXISTS btree_gist;


ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_no_overlap;


ALTER TABLE bookings
ADD COLUMN reserved_until TIMESTAMPTZ NOT NULL;


ALTER TABLE bookings
ADD CONSTRAINT bookings_reserved_until_valid
CHECK (reserved_until >= end_time);


ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tstzrange(start_time, reserved_until) WITH &&
)
WHERE (status = 'booked');


DROP INDEX IF EXISTS idx_bookings_port_time;


CREATE INDEX IF NOT EXISTS idx_bookings_port_start_time
ON bookings(port_id, start_time);

CREATE INDEX IF NOT EXISTS idx_bookings_user_start_time
ON bookings(user_id, start_time DESC);