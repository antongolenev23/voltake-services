DROP INDEX IF EXISTS idx_bookings_user_start_time;


CREATE INDEX idx_bookings_port_time
ON bookings(port_id, start_time, end_time);


ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_no_overlap;


ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_valid_time_range;


DROP EXTENSION IF EXISTS btree_gist;