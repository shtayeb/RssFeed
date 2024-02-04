-- name: StoreNewsLetter :one
INSERT INTO newsletters (created_at, updated_at, email)
VALUES ($1,$2,$3)
RETURNING *;
--

-- name: GetNewsLetter :one
SELECT * FROM newsletters WHERE email = $1;
--

-- name: DeleteNewsLetter :exec
DELETE FROM newsletters WHERE email = $1;  
--

