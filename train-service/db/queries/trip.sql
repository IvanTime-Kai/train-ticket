-- name: CreateTrip :exec
INSERT INTO trips (id, train_id, route_id, departure_time, arrival_time, status)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetTripByID :one
SELECT * FROM trips WHERE id = ? LIMIT 1;

-- name: SearchTrips :many
SELECT t.* FROM trips t
JOIN routes r ON t.route_id = r.id
WHERE r.origin_station_id = ?
AND r.destination_station_id = ?
AND t.departure_time >= ?
AND t.departure_time < ?
AND t.status = 1
ORDER BY t.departure_time ASC;

-- name: UpdateTripStatus :exec
UPDATE trips SET status = ? WHERE id = ?;

-- name: ListTripsByTrain :many
SELECT * FROM trips WHERE train_id = ? ORDER BY departure_time DESC;