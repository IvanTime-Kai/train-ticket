-- name: CreateBooking :exec
INSERT INTO bookings (id, user_id, trip_id, total_price, status, expires_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetBookingByID :one
SELECT * FROM bookings WHERE id = ? LIMIT 1;

-- name: GetBookingsByUserID :many
SELECT * FROM bookings
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: UpdateBookingStatus :exec
UPDATE bookings
SET status = ?, updated_at = NOW()
WHERE id = ?;

-- name: GetExpiredBookings :many
SELECT * FROM bookings
WHERE status = 1
AND expires_at < NOW();

-- name: CreateBookingSeat :exec
INSERT INTO booking_seats (id, booking_id, seat_id, trip_id, seat_number, class, price)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetBookingSeatsByBookingID :many
SELECT * FROM booking_seats
WHERE booking_id = ?;

-- name: GetBookedSeatsByTripID :many
SELECT * FROM booking_seats bs
JOIN bookings b ON bs.booking_id = b.id
WHERE bs.trip_id = ?
AND b.status IN (1, 2);  -- pending hoặc confirmed

-- name: GetBookedSeatsByTripIDForUpdate :many
SELECT bs.seat_id FROM booking_seats bs
JOIN bookings b ON bs.booking_id = b.id
WHERE bs.trip_id = ?
AND b.status IN (1, 2)
FOR UPDATE;

-- name: GetBookingByIDForUpdate :one
SELECT * FROM bookings
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: IsBookingExpired :one
SELECT (expires_at < NOW()) AS expired
FROM bookings
WHERE id = ? LIMIT 1;