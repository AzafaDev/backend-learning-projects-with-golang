-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: CreateUser :one
INSERT INTO users (full_name, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;