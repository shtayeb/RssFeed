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

-- name: GetUserByEmailOrUsername :one
SELECT * FROM users WHERE email=$1 or username=$1;

-- name: UpdateUserName :exec
UPDATE users SET name = $1 WHERE id = $2;

-- name: ChangeUserPassword :exec
UPDATE users SET password = $1 WHERE id = $2;
