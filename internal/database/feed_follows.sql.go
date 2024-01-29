// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"
)

const createFeedFollow = `-- name: CreateFeedFollow :one


INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	FeedID    int32
}

//
func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec

DELETE FROM feed_follows WHERE id = $1 and user_id = $2
`

type DeleteFeedFollowParams struct {
	ID     int32
	UserID int32
}

//
func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.ID, arg.UserID)
	return err
}

const getFeedFollowForUser = `-- name: GetFeedFollowForUser :one

SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type GetFeedFollowForUserParams struct {
	UserID int32
	FeedID int32
}

//
func (q *Queries) GetFeedFollowForUser(ctx context.Context, arg GetFeedFollowForUserParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollowForUser, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_follows.feed_id,feeds.name,feeds.url FROM feed_follows
JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	FeedID    int32
	Name      string
	Url       string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID int32) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.Name,
			&i.Url,
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
