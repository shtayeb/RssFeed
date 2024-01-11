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

// var SessionManager *scs.SessionManager

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
	// SessionManager = scs.New()
	// SessionManager.Lifetime = 24 * time.Second
	// SessionManager.Store = postgresstore.New(db)
	session.InitSessionManager(db)
	authRouter := chi.NewRouter()

	authRouter.Group(func(ar chi.Router) {
		ar.Use(session.SessionManager.LoadAndSave)
		ar.Use(apiCfg.MyMiddleware)
		ar.Get("/home", templ.Handler(views.Home()).ServeHTTP)
		ar.Post("/logout", apiCfg.HandlerLogout)
		ar.Get("/feeds/create", apiCfg.HandlerFeedCreate)
	})

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
	router.Use(session.SessionManager.LoadAndSave)

	router.Get("/", apiCfg.HandlerLandingPage)

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// If user is loggedin redirect them back with a mesasge
		msg := []map[string]string{}
		views.Login(msg).Render(r.Context(), w)
	})
	router.Post("/login", apiCfg.HandlerLogin)

	router.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		// If user is loggedin redirect them back with a mesasge
		msg := []map[string]string{}
		log.Printf("SessionManager-Main: %v", session.SessionManager)
		views.Register(msg, map[string]string{}).Render(r.Context(), w)
	})
	router.Post("/register", apiCfg.HandlerUsersCreate)
	// router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	router.Get("/feeds", apiCfg.HandlerGetFeeds)
	// router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedStore))

	// router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	// router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	// router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	// router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))

	router.Get("/healthz", handlers.HandlerReadiness)
	router.Get("/err", handlers.HandlerErr)

	router.Mount("/", authRouter)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// const collectionConcurrency = 10
	// const collectionInterval = time.Minute
	// go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
