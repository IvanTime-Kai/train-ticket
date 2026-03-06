-- name: CreateRoute :exec
INSERT INTO routes (id, origin_station_id, destination_station_id, distance_km)
VALUES (?, ?, ?, ?);

-- name: GetRouteByID :one
SELECT * FROM routes WHERE id = ? LIMIT 1;

-- name: GetRouteByStations :one
SELECT * FROM routes
WHERE origin_station_id = ?
AND destination_station_id = ?
LIMIT 1;

-- name: ListRoutes :many
SELECT * FROM routes;