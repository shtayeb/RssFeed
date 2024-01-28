package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/shtayeb/rssfeed/http/handlers"
	"github.com/shtayeb/rssfeed/http/session"
	"github.com/shtayeb/rssfeed/internal"
)

type Session struct {
	ID        uuid.UUID
	IPAddress string
	UserAgent string
	Payload   string
	UserID    int
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// log.Println("=========This is for the session redirec: =============")
		// log.Printf("The full url, Host and RequestURI: %v   -----   %v", r.Host, r.RequestURI)
		// log.Printf("This is a try to get origin url: %v", r.Referer())
		// log.Printf("Request reponse - redirect client : %v", r.Response)
		// log.Printf("Request reponse - redirect client : %v", r.UserAgent())
		// log.Println("======================")
		//
		user_id := session.SessionManager.Get(r.Context(), "user_id")

		if user_id == nil {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := internal.DB.GetUserByAPIKey(r.Context(), (user_id).(int32))
		if err != nil {
			// Destroy the sesion
			err := session.SessionManager.Destroy(r.Context())
			if err != nil {
				log.Println("Failed to Destroy the session")
				return
			}

			msgs := []map[string]string{
				{"msg_type": "success", "msg": "Invalid Session, Please login again !"},
			}
			ctx = context.WithValue(r.Context(), "msgs", msgs)
			// redirect back
			// http.Redirect(w, r.WithContext(ctx), "/test", http.StatusSeeOther)
			handlers.RespondWithError(
				w,
				http.StatusNotFound,
				"You are not authenticated - Invalid user_id in the session",
			)
			return
		}
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "user_id", user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTP middleware setting a value on the request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value("user")

		if user == nil {
			// redirect user to login page
			msgs := []map[string]string{
				{"msg_type": "success", "msg": "Please login first"},
			}
			ctx = context.WithValue(ctx, "msgs", msgs)
			// TODO:  context is lost in the redirect
			// r.NewRequestWithContext(ctx)
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
