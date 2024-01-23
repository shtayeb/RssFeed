package handlers

import "github.com/shtayeb/rssfeed/internal/database"

type Config struct {
	AppKey string
	AppEnv string
}

type ApiConfig struct {
	DB     *database.Queries
	Config Config
}
