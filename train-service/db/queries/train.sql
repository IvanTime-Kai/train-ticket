-- name: CreateTrain :exec
INSERT INTO trains (id, name, total_seats, status)
VALUES (?, ?, ?, ?);

-- name: GetTrainByID :one
SELECT * FROM trains WHERE id = ? LIMIT 1;

-- name: ListTrains :many
SELECT * FROM trains WHERE status = 1 ORDER BY name ASC;

-- name: UpdateTrainStatus :exec
UPDATE trains SET status = ? WHERE id = ?;