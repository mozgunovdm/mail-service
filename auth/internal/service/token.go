package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"mts/auth-service/internal/entities/user"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
	Id    string `json:"id"`
}

func GenerateAccessToken(user *user.UserResponse) (string, error) {
	var err error

	claims := &Claims{
		Login: user.Login,
		Id:    user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user *user.UserResponse) (string, error) {
	var err error

	claims := &Claims{
		Login: user.Login,
		Id:    user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseRefreshToken(tokenString string) Claims {
	secret := os.Getenv("REFRESH_SECRET")

	claims := &Claims{}

	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return *claims
}

// func init() {
// 	os.Setenv("ACCESS_SECRET", "123123")  // move to env
// 	os.Setenv("REFRESH_SECRET", "321321") // move to env
// }
