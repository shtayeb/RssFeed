package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func HandlerGetAuthUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user_id").(int32)
	user, err := internal.DB.GetUserByAPIKey(r.Context(), user_id)
	if err != nil {
		// failed to get the user
		log.Println("Faild to get the auth user from user_id in the context")
		// Maybe logout the user
		return
	}

	// Return the template of the user info card

	htmx.NewResponse().
		RenderTempl(r.Context(), w, views.UserInfoCard(user))
}

func HandlerProfileUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(database.User)
	r.ParseForm()
	name := r.PostFormValue("name")

	if name == "" {
		// Invalid data
		// return validation error
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Invalid name sent to server"},
		}
		// ctx := context.WithValue(r.Context(), "msgs", msgs)
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}

	err := internal.DB.UpdateUserName(ctx, database.UpdateUserNameParams{Name: name, ID: user.ID})
	if err != nil {
		// failed to update the DB try again
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Failed to update the user,Please try again !"},
		}
		// ctx := context.WithValue(r.Context(), "msgs", msgs)
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}

	user.Name = name

	msgs := []map[string]string{
		{"msg_type": "success", "msg": "User updated successfully !"},
	}
	ctx = context.WithValue(ctx, "msgs", msgs)
	// return with success/failure message
	htmx.NewResponse().
		AddTrigger(htmx.Trigger("user-updated")).
		RenderTempl(ctx, w, views.UserInfoCard(user))
}

func HanlderUserProfile(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser == nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.UserManagement().Render(r.Context(), w)
}
