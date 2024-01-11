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
	// log.Println(r.Cookies())
	// sessionToken := session.SessionManager.GetString(r.Context(), "session")
	// ctx, err := session.SessionManager.Load(r.Context(), r.Header.Get("X-Session"))
	// log.Printf(
	// 	"User_ID in the handlerLandingPage: %v",
	// 	session.SessionManager.GetString(ctx, "user_id"),
	// )
	//
	// log.Printf("Session Token: %v", sessionToken)

	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	views.Landing(feeds).Render(r.Context(), w)
}
