DROP INDEX IF EXISTS idx_bookings_user_start_time;

DROP INDEX IF EXISTS idx_bookings_port_start_time;


ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_no_overlap;


ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_reserved_until_valid;


ALTER TABLE bookings
DROP COLUMN IF EXISTS reserved_until;


ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tstzrange(start_time, end_time) WITH &&
)
WHERE (status = 'booked');


CREATE INDEX idx_bookings_port_time
ON bookings(port_id, start_time, end_time);