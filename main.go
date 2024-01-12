package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	// "time"

	"github.com/a-h/templ"
	// "github.com/alexedwards/scs/postgresstore"
	// "github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/shtayeb/rssfeed/handlers"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/session"
	"github.com/shtayeb/rssfeed/views"

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
	//
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
	router.Use(middleware.Logger)
	// router.Use(middleware.Recoverer)
	router.Use(session.SessionManager.LoadAndSave)
	router.Use(apiCfg.SessionMiddleware)
	// Public Routes
	router.Group(func(r chi.Router) {
		r.Get("/", apiCfg.HandlerLandingPage)

		r.Get("/login", apiCfg.HandlerLoginView)
		r.Post("/login", apiCfg.HandlerLogin)

		r.Get("/register", apiCfg.HandlerRegisterView)
		r.Post("/register", apiCfg.HandlerUsersCreate)

		r.Get("/feeds", apiCfg.HandlerGetFeeds)

		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerErr)
	})

	// Private Routes
	// Require Authentication

	router.Group(func(ar chi.Router) {
		ar.Use(apiCfg.AuthMiddleware)
		ar.Get("/home", templ.Handler(views.Home()).ServeHTTP)
		ar.Post("/logout", apiCfg.HandlerLogout)

		ar.Get("/feeds", apiCfg.HandlerFeedCreate)
		ar.Post("/feeds", apiCfg.HandlerFeedStore)

		// router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
		// router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
		// router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

		// router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))
		// router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Fecht posts routne
	// const collectionConcurrency = 10
	// const collectionInterval = time.Minute
	// go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
