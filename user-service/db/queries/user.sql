-- name: CreateUser :exec
INSERT INTO users (id, email, password, full_name, phone, role)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET full_name = ?, phone = ?, updated_at = NOW()
WHERE id = ?;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = NOW()
WHERE id = ?;

-- name: VerifyUser :exec
UPDATE users
SET is_verified = 1
WHERE id = ?;

-- name: UpdatePassword :exec
UPDATE users
SET password = ?, updated_at = NOW()
WHERE id = ?;