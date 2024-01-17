package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerLandingPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feeds, err := cfg.DB.GetFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Landing(feeds).Render(ctx, w)
}
