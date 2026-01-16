package infrastructure

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(GetEnv("JWT"))
var JWT_SECRET_REFRESH = []byte(GetEnv("JWT_REFRESH"))

func CreateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}

func CreateRefreshToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET_REFRESH)
}

func RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWT_SECRET_REFRESH, nil
	})

	if err != nil {
		fmt.Println("JWT parse error:", err)
		return "", errors.New("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid refresh token claims")
	}

	if claims["type"] != "refresh" {
		return "", errors.New("Invalid token type")
	}

	userID := uint(claims["user_id"].(float64))
	username := claims["username"].(string)

	return CreateToken(userID, username)
}

func CreateTokenPair(userID uint, username string) (string, string, error) {
	accessToken, err := CreateToken(userID, username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := CreateRefreshToken(userID, username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
