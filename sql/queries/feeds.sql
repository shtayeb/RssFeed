-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeed :one
DELETE FROM feeds WHERE id = $1 AND user_id = $2  
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;
