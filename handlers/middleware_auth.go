package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/session"
)

type Session struct {
	ID        uuid.UUID
	IPAddress string
	UserAgent string
	Payload   string
	UserID    int
}

func refreshToken(cfg ApiConfig, w http.ResponseWriter, s Session) (err error) {
	// newSessionToken := uuid.NewString()

	// Generate new id for the session and update the DB record
	// Delete the older session token

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   (s.ID).String(),
		Expires: time.Now().Add(120 * time.Second),
	})
	return nil
}

// HTTP middleware setting a value on the request context
func (cfg *ApiConfig) MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean UP
		// What if cookie expires on the client and the record exists in our db, It will have unusable data
		// SOLUTION: have a go routine for cleanup
		// session token must to be changed when a user logs in or out of your application

		// ctx, err := sessionpostgresstorer.Context(), r.Header.Get("X-Session"))
		// sessionManager.Put(ctx, "message", "Hello from putMiddleware!")
		ctx := r.Context()
		// sessionToken := session.SessionManager.GetString(ctx, "session_token")

		sessionToken := session.SessionManager.Token(r.Context())
		log.Printf("AuthMiddleware: Session Token : %v", sessionToken)
		// if sessionToken == "" {
		// 	// For any other type of error, return a bad request status
		// 	respondWithError(
		// 		w,
		// 		http.StatusBadRequest,
		// 		"You need to be authenticated ! - sessionToken",
		// 	)
		// 	return
		// }
		//
		// log.Println("sesion_token", sessionToken)
		// 2 - Get the seesion from the DB

		user_id := session.SessionManager.Get(r.Context(), "user_id")

		log.Printf("AuthMiddleware: user_id from session: %v", user_id)
		if user_id == 0 {
			respondWithError(
				w,
				http.StatusBadRequest,
				"You need to be authenticated ! : Bad user_id in session",
			)
			return
		}

		// 3 - Get the user from the BD_Session
		// userID, err := strconv.Atoi(1)
		// if err != nil {
		// 	log.Println("Invalide user_id in the cookie", err)
		// 	respondWithError(w, http.StatusBadRequest, "Invalid user id in the cookie")
		// 	return
		// }
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), (user_id).(int32))
		if err != nil {
			respondWithError(
				w,
				http.StatusNotFound,
				"You are not authenticated - user not found",
			)
			return
		}
		// session.SessionManager.Put(ctx, "user", user)
		// 4 - Set the the user in the Context
		// ctx = context.WithValue(r.Context(), "user", user)
		// 5 - set the session in the context
		ctx = context.WithValue(r.Context(), "user", user)

		// ctx = context.WithValue(r.Context(), "is_authenticated", true)
		// 6 - Move on with the rest of the Request cycle or handler the error redirect

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware here")
		// apiKey, err := auth.GetAPIKey(r.Header)
		// if err != nil {
		// 	respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
		// 	return
		// }
		//
		// user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		// if err != nil {
		// 	respondWithError(w, http.StatusNotFound, "Couldn't get user")
		// 	return
		// }
		// handler(w, r, user)
	}
}
