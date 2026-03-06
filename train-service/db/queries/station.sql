-- name: CreateStation :exec
INSERT INTO stations (id, name, code, city)
VALUES (?, ?, ?, ?);

-- name: GetStationByID :one
SELECT * FROM stations WHERE id = ? LIMIT 1;

-- name: GetStationByCode :one
SELECT * FROM stations WHERE code = ? LIMIT 1;

-- name: ListStations :many
SELECT * FROM stations ORDER BY name ASC;