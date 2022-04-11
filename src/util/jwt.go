package util

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
	"ngb/config"
	"time"
)

type JwtUserClaims struct {
	Id   int  `json:"id"`
	Role bool `json:"role"`
	jwt.StandardClaims
}

var Conf = middleware.JWTConfig{
	Claims:     &JwtUserClaims{},
	SigningKey: []byte(config.C.Jwt.Secret),
}

func GenerateToken(id int, role bool) string {
	claims := &JwtUserClaims{
		id,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.C.Jwt.Secret))
	if err != nil {
		return "error"
	}

	return t
}
