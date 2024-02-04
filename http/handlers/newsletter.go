package handlers

import (
	"log"
	"net/http"
	"time"

	// "github.com/shtayeb/rssfeed/views"

	"github.com/angelofallars/htmx-go"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/views"
)

func StoreNewsLetter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostFormValue("email")

	// Validate the email
	if !rxEmail.Match([]byte(email)) {
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages([]map[string]string{
				{"msg_type": "error", "msg": "Please enter a valid email address"},
			}))
		return
	}

	if email == "" {
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages([]map[string]string{
				{"msg_type": "error", "msg": "You are subscribed successfully!"},
			}))
		return
	}

	// Check if the email already exists
	_, err := internal.DB.GetNewsLetter(r.Context(), email)
	if err == nil {
		log.Println(err)
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages([]map[string]string{
				{"msg_type": "info", "msg": "You are already subscribed !"},
			}))
		return
	}

	_, err = internal.DB.StoreNewsLetter(r.Context(), database.StoreNewsLetterParams{
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Could not add email to newsletter", err)
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages([]map[string]string{
				{"msg_type": "error", "msg": "Could not add email to newsletter !"},
			}))
		return
	}

	// return a HTMX reponse with messages
	htmx.NewResponse().
		RenderTempl(r.Context(), w, views.RenderMessages([]map[string]string{
			{"msg_type": "success", "msg": "You are subscribed successfully!"},
		}))
}
