package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
)

func HandlerFeedFollowsGet(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)

	feedFollows, err := internal.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func HandlerFeedFollowCreate(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)
	type parameters struct {
		FeedID int32
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feedFollow, err := internal.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func HandlerFeedFollowDelete(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := r.Context().Value("user").(database.User)
	feedFollowID, err := strconv.Atoi(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = internal.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     int32(feedFollowID),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct{}{})
}
