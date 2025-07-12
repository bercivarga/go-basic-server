-- name: CreateSession :exec
INSERT INTO sessions (user_id, token, expires_at, refresh_token, refresh_expires_at)
VALUES (?, ?, ?, ?, ?);

-- name: IsValidSession :one
SELECT COUNT(*) FROM sessions
WHERE user_id = ? AND token = ? AND expires_at > CURRENT_TIMESTAMP;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE token = ?;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = ? AND refresh_expires_at > CURRENT_TIMESTAMP;

-- name: DeleteSessionByRefreshToken :exec
DELETE FROM sessions
WHERE refresh_token = ?;
