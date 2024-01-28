package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerLandingPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feeds, err := internal.DB.GetFeeds(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Landing(feeds).Render(ctx, w)
}
