package utils

import (
	"stock-control-back/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JwtSecretKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
