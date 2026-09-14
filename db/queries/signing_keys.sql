-- name: ListSigningKeys :many
-- All signing keys, active first. Feeds the in-memory cache and JWKS.
SELECT *
FROM signing_keys
ORDER BY is_active DESC, created_at DESC;

-- name: GetActiveSigningKey :one
SELECT *
FROM signing_keys
WHERE is_active;

-- name: InsertSigningKey :one
INSERT INTO signing_keys (kid, public_pem, private_enc, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeactivateAllSigningKeys :exec
UPDATE signing_keys SET is_active = false WHERE is_active;

-- name: ActivateSigningKey :execrows
UPDATE signing_keys SET is_active = true WHERE kid = $1;

-- name: DeleteSigningKey :exec
DELETE FROM signing_keys WHERE kid = $1;
