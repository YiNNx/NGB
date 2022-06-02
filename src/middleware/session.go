package middleware

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

type SessionContext struct {
	echo.Context
	Uid  int
	Role bool
}

func HandleSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		if cookie == nil {
			return util.ErrorResponse(c, http.StatusUnauthorized, "haven't logged in yet")
		}
		s, err := model.RedisGet(cookie.Value)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		if s == nil {
			return util.ErrorResponse(c, http.StatusUnauthorized, "authorization out of date")
		}
		session := &Session{}
		err = json.Unmarshal([]byte(s.(string)), session)
		cc := &SessionContext{
			Context: c,
			Uid:     session.Id,
			Role:    session.Role,
		}
		return next(cc)
	}
}

func VerifySuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.(*SessionContext).Role
		if !role {
			err := errors.New("no permission")
			return util.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return next(c)
	}
}

func VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.(*SessionContext).Role
		uid := c.(*SessionContext).Uid

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
			p := &model.Post{Pid: pid}
			err = model.GetByPK(p)
			if err != nil {
				return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
			}
			res, err := model.CheckAdmin(p.Board, uid)
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

type Session struct {
	Id   int  `json:"id"`
	Role bool `json:"role"`
}
