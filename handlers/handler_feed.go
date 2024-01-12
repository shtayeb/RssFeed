package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
	"github.com/shtayeb/rssfeed/views"
)

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
		log.Printf("Failed to create feed", err)
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
