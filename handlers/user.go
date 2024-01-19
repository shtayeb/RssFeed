package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerProfileUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(database.User)
	r.ParseForm()
	name := r.PostFormValue("name")

	if name == "" {
		// Invalid data
		// return validation error
		return
	}

	err := cfg.DB.UpdateUserName(ctx, database.UpdateUserNameParams{Name: name, ID: user.ID})
	if err != nil {
		// failed to update the DB try again
		return
	}

	// return with success/failure message
}

func (cfg *ApiConfig) HanlderUserProfile(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser == nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.UserManagement().Render(r.Context(), w)
}
