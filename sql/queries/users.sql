-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, username,name,email,password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE id = $1;