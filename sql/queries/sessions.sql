-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    provider,
    ip_address,
    user_agent,
    expires_at
)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetSessionByID :one
SELECT s.*, u.*
FROM sessions s
         JOIN users u ON s.user_id = u.id
WHERE s.id = $1 AND s.expires_at > NOW();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();

-- name: GetUserSessions :many
SELECT * FROM sessions
WHERE user_id = $1 AND expires_at > NOW()
ORDER BY created_at DESC;
