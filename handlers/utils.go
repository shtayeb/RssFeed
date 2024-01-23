package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shtayeb/rssfeed/internal/database"
)

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
		return JwtUserInfo{}, fmt.Errorf("Failed to parse the token: %v", err)
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
