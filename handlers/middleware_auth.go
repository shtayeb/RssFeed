package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/shtayeb/rssfeed/internal/session"
)

type Session struct {
	ID        uuid.UUID
	IPAddress string
	UserAgent string
	Payload   string
	UserID    int
}

func (cfg *ApiConfig) SessionMiddleware(next http.Handler) http.Handler {
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

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), (user_id).(int32))
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
			respondWithError(
				w,
				http.StatusNotFound,
				"You are not authenticated - Invalid user_id in the session",
			)
			return
		}
		ctx = context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTP middleware setting a value on the request context
func (cfg *ApiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean UP
		// What if cookie expires on the client and the record exists in our db, It will have unusable data
		// SOLUTION: have a go routine for cleanup
		// session token must to be changed when a user logs in or out of your application

		ctx := r.Context()
		user := ctx.Value("user")

		if user == nil {
			// redirect user to login page

			msgs := []map[string]string{
				{"msg_type": "success", "msg": "Please login first"},
			}
			ctx = context.WithValue(ctx, "msgs", msgs)
			// TODO:  context is lost in the redirect
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
