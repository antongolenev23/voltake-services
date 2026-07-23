ALTER TABLE bookings
DROP CONSTRAINT bookings_no_overlap;

ALTER TABLE bookings
ADD CONSTRAINT bookings_no_overlap
EXCLUDE USING gist (
    port_id WITH =,
    tstzrange(start_time, reserved_until) WITH &&
)
WHERE (status = 'booked');