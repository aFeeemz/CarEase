package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKeyAdmin = []byte(os.Getenv("SECRET_ADMIN"))

type ClaimsAdmin struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWTadmin(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &ClaimsAdmin{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKeyAdmin)
}
