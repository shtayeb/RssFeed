package handlers

import (
	"net/http"
	"strconv"

	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerPostsPage(w http.ResponseWriter, r *http.Request) {
	// templ.Handler(views.NotFoundComponent()).ServeHTTP(w, r)
	ctx := r.Context()
	user := ctx.Value("user").(database.User)
	limitStr := r.URL.Query().Get("limit")
	limit := 9

	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := internal.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}
	views.Posts(posts).Render(ctx, w)
}
