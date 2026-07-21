ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_no_overlap;


ALTER TABLE bookings
    ALTER COLUMN start_time
    TYPE TIMESTAMP
    USING start_time AT TIME ZONE 'UTC';


ALTER TABLE bookings
    ALTER COLUMN end_time
    TYPE TIMESTAMP
    USING end_time AT TIME ZONE 'UTC';


ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tsrange(start_time, end_time) WITH &&
)
WHERE (status = 'booked');

