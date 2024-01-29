-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,feeds.name,feeds.url FROM feed_follows
JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;
--

-- name: GetFeedFollowForUser :one
SELECT * FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
--


-- name: CreateFeedFollow :one
INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
--

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1 and user_id = $2;
--
