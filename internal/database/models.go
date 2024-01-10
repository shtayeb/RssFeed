// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        int32
	LastFetchedAt sql.NullTime
}

type FeedFollow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	FeedID    int32
}

type Post struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      int32
}

type Session struct {
	ID        uuid.UUID
	UserID    int32
	IpAddress sql.NullString
	UserAgent sql.NullString
	Payload   sql.NullString
	ExpiresAt sql.NullTime
}

type User struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Name      string
	Email     string
	Password  string
}
