package handlers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shtayeb/rssfeed/internal/database"
)

type JwtUserInfo struct {
	user_id  int32
	username string
	email    string
	exp      time.Time
}

func createToken(user database.User, cfg ApiConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"user_id":  user.ID,
			"email":    user.Email,
			"exp":      time.Now().Add(time.Minute * 5).Unix(),
		})

	tokenString, err := token.SignedString(cfg.Config.AppKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string, cfg ApiConfig) (JwtUserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	if err != nil {
		return JwtUserInfo{}, fmt.Errorf("Failed to parse the token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return JwtUserInfo{
			user_id:  claims["user_id"].(int32),
			username: claims["username"].(string),
			email:    claims["email"].(string),
			exp:      claims["exp"].(time.Time),
		}, nil
	}

	return JwtUserInfo{}, fmt.Errorf("invalid token")
}
