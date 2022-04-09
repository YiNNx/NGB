package util

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func VerifyUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		idFromToken := c.Get("user").(*jwt.Token).Claims.(*JwtUserClaims).Id
		id, _ := strconv.Atoi(c.Param("id"))
		if id != idFromToken {
			err := errors.New("no permission")
			return ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return next(c)
	}
}

func VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("user").(*jwt.Token).Claims.(*JwtUserClaims).Role
		if !role {
			err := errors.New("no permission")
			return ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return next(c)
	}
}
