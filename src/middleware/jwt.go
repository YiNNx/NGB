package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

func VerifyUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		idFromToken := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
		id, err := strconv.Atoi(c.Param("uid"))
		if err != nil {
			return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		if id != idFromToken {
			err := errors.New("no permission")
			return util.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return next(c)
	}
}

func VerifySuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Role
		if !role {
			err := errors.New("no permission")
			return util.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return next(c)
	}
}

func VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Role
		uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

		if role {
			return next(c)
		}
		if c.Param("bid") != "" {
			bid, err := strconv.Atoi(c.Param("bid"))
			if err != nil {
				return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
			}
			res, err := model.CheckAdmin(bid, uid)
			if err != nil {
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			if !res {
				err := errors.New("no permission")
				return util.ErrorResponse(c, http.StatusForbidden, err.Error())
			}
		}
		if c.Param("pid") != "" {
			pid, err := strconv.Atoi(c.Param("pid"))
			if err != nil {
				return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
			}
			bid, err := model.GetBoardByPost(pid)
			if err != nil {
				return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
			}
			res, err := model.CheckAdmin(bid, uid)
			if err != nil {
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			if !res {
				err := errors.New("no permission")
				return util.ErrorResponse(c, http.StatusForbidden, err.Error())
			}
		}
		return next(c)
	}
}
