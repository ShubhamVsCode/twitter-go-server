package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func GenerateToken(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId.String(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return uuid.MustParse(userId), nil
}
