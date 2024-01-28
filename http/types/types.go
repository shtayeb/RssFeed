package types

type Pagination struct {
	PerPage      int
	CurrentPage  int
	FirstPageUrl string
	LastPageUrl  string
	NextPageUrl  string
	PrevPageUrl  string
	Next         int
	Previous     int
	TotalPage    int
}

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
	APP_KEY string
	APP_ENV string
	APP_URL string
	PORT    string
	MailConfig
}
