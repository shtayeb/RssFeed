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
	res, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Invalid .env file", err)
	}

	APP_KEY := os.Getenv("APP_KEY")
	if APP_KEY == "" {
		log.Fatal("App key is not set. Please create an application key !")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT environment variable is not set")
	}

	DB_URL := os.Getenv("DATABASE_URL")
	if DB_URL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	APP_ENV := os.Getenv("APP_ENV")

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := handlers.ApiConfig{
		DB: dbQueries,
		Config: handlers.Config{
			AppKey:  APP_KEY,
			AppEnv:  APP_ENV,
			APP_URL: res["APP_URL"],
			MailConfig: handlers.MailConfig{
				MAIL_HOST:         res["MAIL_HOST"],
				MAIL_PORT:         res["MAIL_PORT"],
				MAIL_USERNAME:     res["MAIL_USERNAME"],
				MAIL_PASSWORD:     res["MAIL_PASSWORD"],
				MAIL_FROM_ADDRESS: res["MAIL_FROM_ADDRESS"],
				MAIL_FROM_NAME:    res["MAIL_FROM_NAME"],
				MAIL_ENCRYPTION:   res["MAIL_ENCRYPTION"],
			},
		},
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

	if APP_ENV == "" || APP_ENV != "production" {
		router.Use(middleware.Recoverer)
	}

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
		r.Get("/forgot-password", apiCfg.ForgotPasswordView)
		r.Post("/forgot-password", apiCfg.ForgotPassword)

		r.Post("/reset-password", apiCfg.ResetPassword)
		r.Get("/reset-password/{token}", apiCfg.ResetPasswordView)

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
		ar.Get("/feeds/{feedID}/posts", apiCfg.HandlerFeedPosts)
		ar.Post("/feeds", apiCfg.HandlerFeedStore)
		ar.Delete("/feeds/{feedID}", apiCfg.HandlerFeedDelete)

		ar.Get("/feeds/following", apiCfg.HandlerFeedFollowsGet)
		ar.Post("/feeds/following", apiCfg.HandlerFeedFollowCreate)
		ar.Delete("/feeds/following/{feedFollowID}", apiCfg.HandlerFeedFollowDelete)

		ar.Get("/users", apiCfg.HandlerUsersGet)
		ar.Get("/user/profile", apiCfg.HanlderUserProfile)
		ar.Post("/user", apiCfg.HandlerProfileUpdate)
		ar.Get("/user", apiCfg.HandlerGetAuthUser)
		ar.Post("/user/change-password", apiCfg.HandlerChangePassword)
	})

	router.NotFound(apiCfg.HandlerNotFoundPage)

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	// Fecht posts routne
	shouldFetch := true
	if APP_ENV == "" || APP_ENV != "production" {
		log.Println("Not fetching posts right now !")
		shouldFetch = false
	}

	if shouldFetch {
		const collectionConcurrency = 10
		const collectionInterval = time.Minute
		go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
