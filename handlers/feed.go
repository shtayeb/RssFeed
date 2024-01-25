package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerFeedPosts(w http.ResponseWriter, r *http.Request) {
	// Get posts of a specific post
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("size"))
	println(limit)
	if err != nil {
		println("err is nil ")
		limit = 9
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	feed, err := cfg.DB.GetFeed(r.Context(), int32(feedId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	posts, err := cfg.DB.GetFeedPosts(r.Context(), database.GetFeedPostsParams{
		FeedID: int32(feedId),
		Limit:  int32(limit),
		Offset: int32(limit * page),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	totalRecordInDB, _ := cfg.DB.GetFeedPostsCount(r.Context(), feed.ID)
	pagination := paginate(int(totalRecordInDB), limit, page)

	views.FeedPosts(feed, models.DatabaseFeedPostToPostForUserRows(posts), pagination).
		Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerFeedDelete(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	user := r.Context().Value("user").(database.User)
	feed, err := cfg.DB.GetFeed(r.Context(), int32(feedId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feed")
		return
	}

	// // Authorization
	if user.ID != feed.UserID {
		respondWithError(
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
	err = cfg.DB.DeleteFeed(r.Context(), param)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to delete the feed",
		)
		return
	}
	// customize this
	// views.FeedLi(feed).Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerFeedCreate(w http.ResponseWriter, r *http.Request) {
	// feeds/create
	ctx := r.Context()

	feeds, err := cfg.DB.GetFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Feeds(feeds).Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerFeedStore(w http.ResponseWriter, r *http.Request) {
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

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		log.Printf("Failed to create feed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	_, err = cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	views.FeedLi(feed).Render(r.Context(), w)
	// respondWithJSON(w, http.StatusOK, struct {
	// 	feed       models.Feed
	// 	feedFollow models.FeedFollow
	// }{
	// 	feed:       models.DatabaseFeedToFeed(feed),
	// 	feedFollow: models.DatabaseFeedFollowToFeedFollow(feedFollow),
	// })
}

func (cfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
}
