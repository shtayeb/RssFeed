package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerToggleFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.Atoi(chi.URLParam(r, "feedID"))
	if err != nil {
		log.Printf("Failed to parse feedID: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Invalid feed ID")
		return
	}

	msgs := []map[string]string{}
	// Get feed follow
	user := r.Context().Value("user").(database.User)
	feedFollow, err := internal.DB.GetFeedFollowForUser(
		r.Context(),
		database.GetFeedFollowForUserParams{
			UserID: int32(user.ID),
			FeedID: int32(feedId),
		},
	)
	if err != nil {
		// follow
		log.Printf("Follow: %v", err)
		user := r.Context().Value("user").(database.User)
		_, err := internal.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    int32(feedId),
		})
		if err != nil {
			msgs = []map[string]string{
				{
					"msg_type": "error",
					"msg":      "Could'nt follow the feed, try again!",
				},
			}
		} else {
			msgs = []map[string]string{
				{
					"msg_type": "success",
					"msg":      "Feed has been followed",
				},
			}
		}
	} else {
		log.Printf("Unfollow")
		// Unfollow
		err = internal.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
			UserID: user.ID,
			ID:     int32(feedFollow.ID),
		})

		if err != nil {
			msgs = []map[string]string{
				{
					"msg_type": "error",
					"msg":      "Could'nt unfollow the feed, try again!",
				},
			}
		} else {
			msgs = []map[string]string{
				{
					"msg_type": "info",
					"msg":      "Feed has been Unfollowed",
				},
			}
		}

	}

	htmx.NewResponse().
		RenderTempl(r.Context(), w, views.RenderMessages(msgs))
}
