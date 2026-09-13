-- name: InsertRefreshToken :one
-- Insert a new refresh token for a user.
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, ip_address)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRefreshToken :one
-- Retrieve a refresh token by its hash.
SELECT *
FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshToken :exec
-- Delete a refresh token by its hash.
DELETE FROM refresh_tokens
WHERE token_hash = $1;
