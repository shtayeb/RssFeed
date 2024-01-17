package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/shtayeb/rssfeed/handlers"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/session"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := handlers.ApiConfig{
		DB: dbQueries,
	}

	// Sesison Manager
	session.InitSessionManager(db)

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
	// router.Use(middleware.Recoverer)
	router.Use(session.SessionManager.LoadAndSave)
	router.Use(apiCfg.SessionMiddleware)

	// Static files handler
	fs := http.FileServer(http.Dir("public"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public Routes
	router.Group(func(r chi.Router) {
		r.Get("/", apiCfg.HandlerLandingPage)
		r.Get("/login", apiCfg.HandlerLoginView)
		r.Post("/login", apiCfg.HandlerLogin)

		r.Get("/register", apiCfg.HandlerRegisterView)
		r.Post("/register", apiCfg.HandlerUsersCreate)

		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerErr)
	})

	// Private Routes - Require Authentication
	router.Group(func(ar chi.Router) {
		ar.Use(apiCfg.AuthMiddleware)
		ar.Get("/posts", apiCfg.HandlerPostsPage)
		ar.Post("/logout", apiCfg.HandlerLogout)

		ar.Get("/feeds", apiCfg.HandlerFeedCreate)
		ar.Post("/feeds", apiCfg.HandlerFeedStore)
		ar.Delete("/feeds/{feedID}", apiCfg.HandlerFeedDelete)

		ar.Get("/feeds/following", apiCfg.HandlerFeedFollowsGet)
		ar.Post("/feeds/following", apiCfg.HandlerFeedFollowCreate)
		ar.Delete("/feeds/following/{feedFollowID}", apiCfg.HandlerFeedFollowDelete)

		ar.Get("/users", apiCfg.HandlerUsersGet)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Fecht posts routne
	appEnv := os.Getenv("APP_ENV")
	shouldFetch := true
	if appEnv == "" || appEnv != "production" {
		log.Fatal("Not fetching posts right now !")
		shouldFetch = false
	}

	if shouldFetch {
		const collectionConcurrency = 10
		const collectionInterval = time.Minute
		go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
