// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    int32
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const deleteFeed = `-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1 AND user_id = $2
`

type DeleteFeedParams struct {
	ID     int32
	UserID int32
}

func (q *Queries) DeleteFeed(ctx context.Context, arg DeleteFeedParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeed, arg.ID, arg.UserID)
	return err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE id = $1
`

func (q *Queries) GetFeed(ctx context.Context, id int32) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedPosts = `-- name: GetFeedPosts :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id, feeds.name as feed_name,feeds.url as feed_url 
FROM posts 
JOIN feeds ON feeds.id = posts.feed_id
WHERE feed_id = $1
ORDER BY posts.published_at DESC
LIMIT $2 offset $3
`

type GetFeedPostsParams struct {
	FeedID int32
	Limit  int32
	Offset int32
}

type GetFeedPostsRow struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      int32
	FeedName    string
	FeedUrl     string
}

func (q *Queries) GetFeedPosts(ctx context.Context, arg GetFeedPostsParams) ([]GetFeedPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedPosts, arg.FeedID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedPostsRow
	for rows.Next() {
		var i GetFeedPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.FeedName,
			&i.FeedUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedPostsCount = `-- name: GetFeedPostsCount :one
SELECT COUNT(*) FROM posts WHERE feed_id = $1
`

func (q *Queries) GetFeedPostsCount(ctx context.Context, feedID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getFeedPostsCount, feedID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
ORDER BY feeds.created_at DESC
LIMIT $1 offset $2
`

type GetFeedsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetFeeds(ctx context.Context, arg GetFeedsParams) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedsCount = `-- name: GetFeedsCount :one
SELECT COUNT(*) FROM feeds
`

func (q *Queries) GetFeedsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getFeedsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getNextFeedsToFetch = `-- name: GetNextFeedsToFetch :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1
`

func (q *Queries) GetNextFeedsToFetch(ctx context.Context, limit int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextFeedsToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markFeedFetched = `-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id int32) (Feed, error) {
	row := q.db.QueryRowContext(ctx, markFeedFetched, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}
