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

-- name: GetUserByEmail :one
-- Retrieves a user by their unique email address.
SELECT *
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
-- Inserts a new user into the users table and returns the created user.
INSERT INTO users (username, email, password_hash, roles)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
-- Updates an existing user's information
UPDATE users
SET username = COALESCE(sqlc.narg('username'), username),
    email = COALESCE(sqlc.narg('email'), email),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    roles = COALESCE(sqlc.narg('roles'), roles),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
-- Deletes a user from the users table by their unique identifier.
DELETE FROM users
WHERE id = $1;

-- name: UsernameExists :one
-- Checks if a username already exists in the users table.
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE username = $1
);

-- name: EmailExists :one
-- Checks if an email address already exists in the users table.
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
);

-- name: ListUsers :many
-- Retrieves a list of users with optional pagination.
SELECT username, email, roles, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
