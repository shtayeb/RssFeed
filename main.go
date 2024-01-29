package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shtayeb/rssfeed/http/handlers"
	"github.com/shtayeb/rssfeed/http/middlewares"
	"github.com/shtayeb/rssfeed/http/session"
	"github.com/shtayeb/rssfeed/internal"

	_ "github.com/lib/pq"
)

func main() {
	// Loads ENV variables, initializes DB and Session Manager
	internal.InitApp()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Middlewares
	router.Use(middleware.Logger)

	if internal.Config.APP_ENV == "" || internal.Config.APP_ENV != "production" {
		router.Use(middleware.Recoverer)
	}

	router.Use(session.SessionManager.LoadAndSave)
	router.Use(middlewares.SessionMiddleware)

	// Static files handler
	fs := http.FileServer(http.Dir("public"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public Routes
	router.Group(func(r chi.Router) {
		r.Get("/", handlers.HandlerLandingPage)
		r.Get("/login", handlers.HandlerLoginView)
		r.Post("/login", handlers.HandlerLogin)
		r.Get("/forgot-password", handlers.ForgotPasswordView)
		r.Post("/forgot-password", handlers.ForgotPassword)

		r.Post("/reset-password", handlers.ResetPassword)
		r.Get("/reset-password/{token}", handlers.ResetPasswordView)

		r.Get("/register", handlers.HandlerRegisterView)
		r.Post("/register", handlers.HandlerUsersCreate)

		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerErr)
	})

	// Private Routes - Require Authentication
	router.Group(func(ar chi.Router) {
		ar.Use(middlewares.AuthMiddleware)
		ar.Get("/posts", handlers.HandlerPostsPage)
		ar.Post("/logout", handlers.HandlerLogout)

		ar.Get("/feeds", handlers.HandlerGetFeeds)
		ar.Get("/user/feeds", handlers.HandlerFeedCreate)
		ar.Get("/feeds/{feedID}/posts", handlers.HandlerFeedPosts)
		ar.Post("/feeds", handlers.HandlerFeedStore)
		ar.Delete("/feeds/{feedID}", handlers.HandlerFeedDelete)

		// ar.Get("/feeds/following", handlers.HandlerFeedFollowsGet)
		// ar.Post("/feeds/following", handlers.HandlerFeedFollowCreate)
		ar.Put("/feeds/following/{feedID}", handlers.HandlerToggleFeedFollow)

		ar.Get("/users", handlers.HandlerUsersGet)
		ar.Get("/user/profile", handlers.HanlderUserProfile)
		ar.Post("/user", handlers.HandlerProfileUpdate)
		ar.Get("/user", handlers.HandlerGetAuthUser)
		ar.Post("/user/change-password", handlers.HandlerChangePassword)
	})

	router.NotFound(handlers.HandlerNotFoundPage)

	srv := &http.Server{
		Addr:    ":" + internal.Config.PORT,
		Handler: router,
	}

	// Fecht posts routne
	shouldFetch := true
	if internal.Config.APP_ENV == "" || internal.Config.APP_ENV != "production" {
		log.Println("Not fetching posts right now !")
		shouldFetch = false
	}

	if shouldFetch {
		const collectionConcurrency = 10
		const collectionInterval = time.Minute
		go internal.StartScraping(internal.DB, collectionConcurrency, collectionInterval)
	}

	log.Printf("Serving on port: %s\n", internal.Config.PORT)
	log.Fatal(srv.ListenAndServe())
}
