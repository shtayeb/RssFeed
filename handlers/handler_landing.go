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
	// templ.Handler(views.NotFoundComponent()).ServeHTTP(w, r)
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Landing(feeds).Render(r.Context(), w)
}
