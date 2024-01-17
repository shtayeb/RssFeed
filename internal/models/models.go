package models

import (
	"database/sql"
	"time"

	"github.com/shtayeb/rssfeed/internal/database"
)

type User struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

func DatabaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
	}
}

type Feed struct {
	ID            int32      `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        int32      `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func DatabaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: nullTimeToTimePtr(feed.LastFetchedAt),
	}
}

func DatabaseFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, feed := range feeds {
		result[i] = DatabaseFeedToFeed(feed)
	}
	return result
}

type FeedFollow struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int32     `json:"user_id"`
	FeedID    int32     `json:"feed_id"`
}

func DatabaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}
}

func DatabaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollow {
	result := make([]FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		result[i] = DatabaseFeedFollowToFeedFollow(feedFollow)
	}
	return result
}

type Post struct {
	ID          int32      `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      int32      `json:"feed_id"`
}

func DatabaseFeedPostToPostForUserRow(post database.GetFeedPostsRow) database.GetPostsForUserRow {
	return database.GetPostsForUserRow(post)
}

func DatabaseFeedPostToPostForUserRows(posts []database.GetFeedPostsRow) []database.GetPostsForUserRow {
	result := make([]database.GetPostsForUserRow, len(posts))
	for i, post := range posts {
		result[i] = DatabaseFeedPostToPostForUserRow(post)
	}
	return result
}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

// func nullStringToStringPtr(s sql.NullString) *string {
// 	if s.Valid {
// 		return &s.String
// 	}
// 	return nil
// }
