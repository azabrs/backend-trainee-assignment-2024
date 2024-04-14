package jwt

import (
	custom_errors "backend-trainee-assignment-2024/errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	*jwt.RegisteredClaims
	IsAdmin bool `json:"isAdmin"`
}

func CreateJWT(isAdmin bool, signingKey []byte, expiresAt time.Time) (string, error) {
	claims := &Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		IsAdmin: isAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckIsAdminInJWT(tokenString string, signingKey string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil{
		return false, err
	}
	if !token.Valid{
		return false, custom_errors.ErrAdminRequired
	}

	return claims.IsAdmin, nil
}