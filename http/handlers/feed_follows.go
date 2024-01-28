package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

// GET: feeds/following
// List all feeds that a user is following
func HandlerFeedFollowsGet(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)

	feedFollows, err := internal.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could'nt get your following feeds")
	}

	RespondWithJSON(w, http.StatusOK, feedFollows)
}

func HandlerFeedFollowCreate(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	user := r.Context().Value("user").(database.User)

	feedFollow, err := internal.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    int32(feedId),
	})

	log.Printf("feedFollow: %v", feedFollow)
	if err != nil {
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Couldn't create feed follow"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		// RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
	}

	msgs := []map[string]string{
		{"msg_type": "success", "msg": "You have followed the feed"},
	}
	htmx.NewResponse().
		RenderTempl(r.Context(), w, views.RenderMessages(msgs))

	// RespondWithJSON(w, http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func HandlerFeedFollowDelete(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := r.Context().Value("user").(database.User)
	feedFollowID, err := strconv.Atoi(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Couldn't decode parameters"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))

		return
		// RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	err = internal.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     int32(feedFollowID),
	})
	if err != nil {
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Couldn't create feed follow"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		// RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}
	msgs := []map[string]string{
		{"msg_type": "success", "msg": "Couldn't create feed follow"},
	}
	htmx.NewResponse().
		RenderTempl(r.Context(), w, views.RenderMessages(msgs))

	// RespondWithJSON(w, http.StatusOK, struct{}{})
}
