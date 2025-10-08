package util

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJsonWebToken(userID uuid.UUID) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    "go_social",
		Subject:   userID.String(),
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

func HandleBindError(c *gin.Context, err error) {

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(400, gin.H{"error": "invalid input"})
	}

	result := make(map[string]string)

	for _, fe := range validationErrors {
		errMsg := GetErrorMessage(fe)
		result[fe.Field()] = errMsg

	}
	c.JSON(400, gin.H{"error": result})
}

func GetErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "alphanum":
		return "should be alpanumeric"
	case "uuid":
		return "should be uuid"
	case "email":
		return "should be valid email"
	case "lte":
		return fmt.Sprintf("should less than or equal to %s", fe.Param())
	}
	return "not valid"
}
