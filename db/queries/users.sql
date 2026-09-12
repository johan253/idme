-- name: GetUserById :one
-- Retrieves a user by their unique identifier.
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
-- Retrieves a user by their unique username.
SELECT *
FROM users
WHERE username = $1 LIMIT 1;
