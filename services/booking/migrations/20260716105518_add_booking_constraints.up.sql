CREATE EXTENSION IF NOT EXISTS btree_gist;


ALTER TABLE bookings
ADD CONSTRAINT bookings_valid_time_range
CHECK (end_time > start_time);


ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tsrange(start_time, end_time) WITH &&
)
WHERE (status = 'booked');

DROP INDEX IF EXISTS idx_bookings_port_time;


CREATE INDEX idx_bookings_user_start_time
ON bookings(user_id, start_time DESC);