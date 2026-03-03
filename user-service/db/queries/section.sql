-- name: CreateSession :exec
INSERT INTO user_sessions (id, user_id, device, ip_address, logged_in_at)
VALUES (?, ?, ?, ?, NOW());

-- name: GetSessionsByUserID :many
SELECT * FROM user_sessions
WHERE user_id = ?
ORDER BY logged_in_at DESC;

-- name: UpdateSessionLogout :exec
UPDATE user_sessions
SET logged_out_at = NOW()
WHERE id = ?;

-- name: LogoutAllSessions :exec
UPDATE user_sessions
SET logged_out_at = NOW()
WHERE user_id = ? AND logged_out_at IS NULL;