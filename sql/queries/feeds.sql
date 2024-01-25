-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1 AND user_id = $2;  

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

-- name: GetFeedPosts :many
SELECT posts.*, feeds.name as feed_name,feeds.url as feed_url 
FROM posts 
JOIN feeds ON feeds.id = posts.feed_id
WHERE feed_id = $1
ORDER BY posts.published_at DESC
LIMIT $2 offset $3;

-- name: GetFeedPostsCount :one
SELECT COUNT(*) FROM posts WHERE feed_id = $1;

