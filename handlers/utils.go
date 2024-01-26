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
	"github.com/shtayeb/rssfeed/internal/types"
)

func paginate(totalData int, limit int, page int) types.Pagination {
	paginated := types.Pagination{}
	// get count of the all of the tables data

	// Count all record
	totalPage := (totalData / limit)
	// log.Printf("totalPage = data(%v)/limit(%v): %v", totalData, limit, totalPage)

	// Calculator Total Page
	remainder := (totalData % limit)
	log.Printf("Remainder totalData limit = %v", remainder)

	if remainder == 0 {
		paginated.TotalPage = totalPage
	} else {
		paginated.TotalPage = totalPage + 1
	}

	// log.Printf("======== Pagination.TotalPage : %v  ========= ", paginated.TotalPage)
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
		paginated.Next = 1
	}

	paginated.FirstPageUrl = fmt.Sprintf("?page=%v&size=%v", 1, limit)
	paginated.LastPageUrl = fmt.Sprintf("?page=%v&size=%v", paginated.TotalPage, limit)
	paginated.NextPageUrl = fmt.Sprintf("?page=%v&size=%v", paginated.Next, limit)
	paginated.PrevPageUrl = fmt.Sprintf("?page=%v&size=%v", paginated.Previous, limit)

	// log.Printf("==== This is inside the paginate === : %v \n", paginated)
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

// // some generic typing here
// func paginatePlus(
// 	getData func(context.Context, any) ([]string, error),
// 	r http.Request,
// ) types.Pagination {
// 	// More generic and full fucntion paginate
// 	paginated := types.Pagination{}
// 	limit, err := strconv.Atoi(r.URL.Query().Get("size"))
// 	if err != nil {
// 		limit = 9
// 	}
// 	page, err := strconv.Atoi(r.URL.Query().Get("page"))
// 	if err != nil || page <= 0 {
// 		page = 1
// 	}
//
// 	// Get the data
// 	params := "test"
// 	// params := {
// 	// 	FeedID: int32(feedId),
// 	// 	Limit:  int32(limit),
// 	// 	Offset: int32(limit * (page - 1)),
// 	// }
// 	log.Println(limit)
// 	data, err := getData(r.Context(), params)
// 	if err != nil {
// 		log.Println(data)
// 		log.Println(err)
// 	}
//
// 	return paginated
// }
