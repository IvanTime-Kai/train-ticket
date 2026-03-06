-- name: CreateSeat :exec
INSERT INTO seats (id, train_id, seat_number, class, price)
VALUES (?, ?, ?, ?, ?);

-- name: GetSeatByID :one
SELECT * FROM seats WHERE id = ? LIMIT 1;

-- name: ListSeatsByTrain :many
SELECT * FROM seats WHERE train_id = ? ORDER BY seat_number ASC;