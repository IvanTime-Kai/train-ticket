-- name: CreatePayment :exec
INSERT INTO payments (id, booking_id, amount, method, status)
VALUES (?, ?, ?, ?, ?);

-- name: GetPaymentByBookingID :one
SELECT * FROM payments
WHERE booking_id = ? LIMIT 1;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = ?, paid_at = NOW()
WHERE booking_id = ?;