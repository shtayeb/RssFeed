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

// func refreshToken(cfg ApiConfig, w http.ResponseWriter, s Session) (err error) {
// 	// newSessionToken := uuid.NewString()
//
// 	// Generate new id for the session and update the DB record
// 	// Delete the older session token
//
// 	// Set the new token as the users `session_token` cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Value:   (s.ID).String(),
// 		Expires: time.Now().Add(120 * time.Second),
// 	})
// 	return nil
// }

func (cfg *ApiConfig) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
			// redirect user back, with message
			respondWithError(
				w,
				http.StatusBadRequest,
				"You need to be authenticated, No user from session",
			)
			return
		}

		// you will have user_id in the context from the sesison middleware
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
