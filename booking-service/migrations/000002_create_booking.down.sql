DROP INDEX IF EXISTS idx_status_history_booking;
DROP INDEX IF EXISTS idx_bookings_expires_at;
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_car_id;
DROP INDEX IF EXISTS idx_bookings_user_id;

DROP TABLE IF EXISTS booking_status_history;
DROP TABLE IF EXISTS bookings;
