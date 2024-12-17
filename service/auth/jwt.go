package auth

import (
	"strconv"
	"time"

	"github.com/AdvenAdam/go-ecom/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(secret []byte, userID int) (string, error) {
	Expire := time.Second * time.Duration(config.Envs.JWTExpireInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(Expire),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
