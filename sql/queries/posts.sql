-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
--

-- name: GetPostsForUser :many
SELECT posts.*,feeds.name as feed_name,feeds.url as feed_url FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN feeds on feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
 --
