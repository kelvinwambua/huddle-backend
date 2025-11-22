-- name: CreateOAuthUser :one
INSERT INTO users (
    username,
    email,
    avatar_url,
    provider,
    provider_user_id,
    access_token,
    refresh_token,
    expires_at,
    name,
    first_name,
    last_name,
    nick_name,
    description,
    location
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    RETURNING *;

-- name: GetUserByProviderID :one
SELECT * FROM users
WHERE provider = $1 AND provider_user_id = $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserOAuthTokens :one
UPDATE users
SET
    access_token = $2,
    refresh_token = $3,
    expires_at = $4,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2, avatar_url = $3, updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
