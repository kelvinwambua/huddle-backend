-- name: CreateProfile :one
INSERT INTO profiles (
    user_id,
    username,
    display_name,
    bio,
    website
) VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetProfileByUserID :one
SELECT * FROM profiles
WHERE user_id = $1;

-- name: GetProfileByUsername :one
SELECT * FROM profiles
WHERE username = $1;

-- name: CheckUsernameExists :one
SELECT EXISTS(
    SELECT 1 FROM profiles
    WHERE username = $1
) AS exists;

-- name: UpdateProfile :one
UPDATE profiles
SET
    username = $2,
    display_name = $3,
    bio = $4,
    website = $5,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: UpdateUsername :one
UPDATE profiles
SET
    username = $2,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE user_id = $1;

-- name: ListProfiles :many
SELECT * FROM profiles
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: SearchProfilesByUsername :many
SELECT * FROM profiles
WHERE username ILIKE $1
ORDER BY username
    LIMIT $2 OFFSET $3;

