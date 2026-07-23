CREATE TYPE booking_status_new AS ENUM (
    'booked',
    'cancelled',
    'completed'
);

ALTER TABLE bookings
    ALTER COLUMN status TYPE booking_status_new
    USING status::text::booking_status_new;

DROP TYPE booking_status;

ALTER TYPE booking_status_new RENAME TO booking_status;
