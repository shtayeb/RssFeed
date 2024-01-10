-- name: CreateSession :one
INSERT INTO sessions (id,user_id,ip_address,user_agent,payload,expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE id = $1;


