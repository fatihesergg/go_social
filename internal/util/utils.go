package util

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJsonWebToken(userID int64) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    "go_social",
		Subject:   strconv.FormatInt(userID, 10),
		Audience:  jwt.ClaimStrings{"go_social_user"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),

		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	secret := os.Getenv("JWT_SECRET")

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
