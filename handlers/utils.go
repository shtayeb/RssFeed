package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shtayeb/rssfeed/internal/database"
)

type Pagination struct {
	PerPage      int
	CurrentPage  int
	LastPage     int
	FirstPageUrl string
	LastPageUrl  string
	NextPageUrl  string
	PrevPageUrl  string
	Next         int
	Previous     int
	TotalPage    int
	Data         *[]any
}

// Generated Pagination Meta data
func paginate[T interface{}](data []T, limit int, page int) Pagination {
	paginated := Pagination{}

	// Count all record
	total := (len(data) / limit)

	// Calculator Total Page
	remainder := (total % limit)
	if remainder == 0 {
		paginated.TotalPage = total
	} else {
		paginated.TotalPage = total + 1
	}

	// Set current/record per page meta data
	paginated.CurrentPage = page
	paginated.PerPage = limit

	// Calculator the Next/Previous Page
	if page <= 0 {
		paginated.Next = page + 1
	} else if page < paginated.TotalPage {
		paginated.Previous = page - 1
		paginated.Next = page + 1
	} else if page == paginated.TotalPage {
		paginated.Previous = page - 1
		paginated.Next = 0
	}

	return paginated
}

type MailRequest struct {
	to      []string
	subject string
}

func SendEmail(
	c context.Context,
	comp templ.Component,
	appConfig Config,
	r MailRequest,
) (bool, error) {
	auth := smtp.PlainAuth(
		"",
		appConfig.MAIL_USERNAME,
		appConfig.MAIL_PASSWORD,
		appConfig.MAIL_HOST,
	)

	// parse template and set it to r.Body
	var mailBody bytes.Buffer
	err := comp.Render(c, &mailBody)
	if err != nil {
		return false, fmt.Errorf("failed to parse the email template: %v", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + mailBody.String())

	if err := smtp.SendMail(appConfig.MAIL_HOST+":"+appConfig.MAIL_PORT, auth, appConfig.MAIL_FROM_ADDRESS, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func RenderWithMsg(
	component templ.Component,
	w http.ResponseWriter,
	c context.Context,
	msgs []map[string]string,
) error {
	ctx := context.WithValue(c, "msgs", msgs)
	return component.Render(ctx, w)
}

func Render(
	component templ.Component,
	w http.ResponseWriter,
	c context.Context,
) error {
	return component.Render(c, w)
}

type JwtUserInfo struct {
	user_id  int32
	username string
	email    string
	exp      int64
}

func createToken(user database.User, cfg ApiConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"user_id":  user.ID,
			"email":    user.Email,
			"exp":      time.Now().Add(time.Minute * 5).Unix(),
		})

	log.Printf("appkey in the createToken: %v", cfg.Config.AppKey)

	tokenString, err := token.SignedString([]byte(cfg.Config.AppKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string, cfg ApiConfig) (JwtUserInfo, error) {
	// log.Printf("token to verifyToken: %v", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// log.Printf("Token inside keyfunc: %v", token.Claims.(jwt.MapClaims)["email"])
		return []byte(cfg.Config.AppKey), nil
	})
	if err != nil {
		return JwtUserInfo{}, fmt.Errorf("failed to parse the token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	log.Printf("========== token after Claims =======: %v", claims)
	if ok && token.Valid {
		return JwtUserInfo{
			user_id:  int32(claims["user_id"].(float64)),
			username: claims["username"].(string),
			email:    claims["email"].(string),
			exp:      int64(claims["exp"].(float64)),
		}, nil
	}

	return JwtUserInfo{}, fmt.Errorf("invalid token")
}
