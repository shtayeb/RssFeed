package handlers

import (
	"net/http"
	"strconv"

	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerPostsPage(w http.ResponseWriter, r *http.Request) {
	// templ.Handler(views.NotFoundComponent()).ServeHTTP(w, r)
	ctx := r.Context()
	user := ctx.Value("user").(database.User)
	limitStr := r.URL.Query().Get("limit")
	limit := 9

	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}
	views.Posts(posts).Render(ctx, w)
}

// func (cfg *ApiConfig) HandlerPostsGet(w http.ResponseWriter, r *http.Request) {
// user := r.Context().Value("user").(database.User)
// limitStr := r.URL.Query().Get("limit")
// limit := 10
// if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
// 	limit = specifiedLimit
// }
//
// posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
// 	UserID: user.ID,
// 	Limit:  int32(limit),
// })
// if err != nil {
// 	respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
// 	return
// }
// return
// // respondWithJSON(w, http.StatusOK, models.DatabasePostsToPosts(posts))
// }
