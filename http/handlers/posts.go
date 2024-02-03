package handlers

import (
	"net/http"
	"strconv"

	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerPostsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	limit, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		println("err is nil ")
		limit = 9
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	posts, err := internal.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: int32(limit * (page - 1)),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}

	totalRecordInDB, _ := internal.DB.GetPostForUserCount(r.Context(), user.ID)
	pagination := paginate(int(totalRecordInDB), limit, page)

	views.Posts(posts, pagination).Render(ctx, w)
}
