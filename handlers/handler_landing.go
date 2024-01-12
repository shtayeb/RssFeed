package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerTestPage(w http.ResponseWriter, r *http.Request) {
	// templ.Handler(views.NotFoundComponent()).ServeHTTP(w, r)
	respondWithError(w, http.StatusInternalServerError, "success here in the middleware")
	return
}

func (cfg *ApiConfig) HandlerLandingPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feeds, err := cfg.DB.GetFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Landing(feeds).Render(ctx, w)
}
