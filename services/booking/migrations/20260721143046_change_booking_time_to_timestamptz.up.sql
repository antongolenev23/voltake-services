
ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_no_overlap;


ALTER TABLE bookings
    ALTER COLUMN start_time
    TYPE TIMESTAMPTZ
    USING start_time AT TIME ZONE 'UTC';


ALTER TABLE bookings
    ALTER COLUMN end_time
    TYPE TIMESTAMPTZ
    USING end_time AT TIME ZONE 'UTC';


ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tstzrange(start_time, end_time) WITH &&
)
WHERE (status = 'booked');
