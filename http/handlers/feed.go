package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerFeedPosts(w http.ResponseWriter, r *http.Request) {
	// Get posts of a specific post
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		println("err is nil ")
		limit = 9
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	feed, err := internal.DB.GetFeed(r.Context(), int32(feedId))
	if err != nil {
		log.Printf("Here is why: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	posts, err := internal.DB.GetFeedPosts(r.Context(), database.GetFeedPostsParams{
		FeedID: int32(feedId),
		Limit:  int32(limit),
		Offset: int32(limit * (page - 1)),
	})
	if err != nil {
		log.Printf("Here is why: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	totalRecordInDB, _ := internal.DB.GetFeedPostsCount(r.Context(), feed.ID)
	pagination := paginate(int(totalRecordInDB), limit, page)

	views.FeedPosts(feed, models.DatabaseFeedPostToPostForUserRows(posts), pagination).
		Render(r.Context(), w)
}

func HandlerFeedDelete(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	user := r.Context().Value("user").(database.User)
	feed, err := internal.DB.GetFeed(r.Context(), int32(feedId))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	// // Authorization
	if user.ID != feed.UserID {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"You are not authorized to delete this feed",
		)
		return
	}

	param := database.DeleteFeedParams{
		ID:     feed.ID,
		UserID: user.ID,
	}
	err = internal.DB.DeleteFeed(r.Context(), param)
	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to delete the feed",
		)
		return
	}
	// customize this
	// views.FeedLi(feed).Render(r.Context(), w)
}

func HandlerFeedCreate(w http.ResponseWriter, r *http.Request) {
	// feeds/create
	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	// feeds, err := internal.DB.GetFeeds(ctx)
	feedFollows, err := internal.DB.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	log.Println("YOOOOOOOOOOOOOOOOO", feedFollows)

	views.Feeds(feedFollows).Render(ctx, w)
}

func HandlerFeedStore(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	user := r.Context().Value("user").(database.User)

	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	params := parameters{
		Name: r.PostFormValue("name"),
		URL:  r.PostFormValue("url"),
	}

	log.Printf("FeedName: %v ---- FeedURL: %v", params.Name, params.URL)

	feed, err := internal.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		log.Printf("Failed to create feed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	feedFollow, err := internal.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	views.FeedLi(database.GetFeedFollowsForUserRow{
		Name:      feed.Name,
		Url:       feed.Url,
		FeedID:    feed.ID,
		UserID:    feedFollow.UserID,
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}).Render(r.Context(), w)
}

func HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := internal.DB.GetFeeds(r.Context())
	if err != nil {
		log.Printf("failed to get feeds from DB: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		limit = 9
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	totalRecordInDB, _ := internal.DB.GetFeedsCount(r.Context())
	pagination := paginate(int(totalRecordInDB), limit, page)

	views.AllFeeds(feeds, pagination).Render(r.Context(), w)
}
