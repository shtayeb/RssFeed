package handlers

import "github.com/shtayeb/rssfeed/internal/database"

type MailConfig struct {
	MAIL_HOST         string
	MAIL_PORT         string
	MAIL_USERNAME     string
	MAIL_PASSWORD     string
	MAIL_ENCRYPTION   string
	MAIL_FROM_ADDRESS string
	MAIL_FROM_NAME    string
}

type Config struct {
	AppKey string
	AppEnv string
	MailConfig
}

type ApiConfig struct {
	DB     *database.Queries
	Config Config
}
