package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	key []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{
		key: []byte(secret),
	}
}

func (j *JWT) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.key)
}

func (j *JWT) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.key, nil
	})
}

func (j *JWT) ExtractUserId(token *jwt.Token) (int, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["user_id"].(float64); ok {
			return int(userId), nil
		}
	}
	return 0, fmt.Errorf("invalid token")
}
