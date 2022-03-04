package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type customClaims struct {
	Username string `json:"username"`
	Level    int    `json:"level"`
	jwt.StandardClaims
}

func ParseToken(tokenStr string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("hello"), nil
	})
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
