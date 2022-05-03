package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	"ngb/util"
)

func initMessageRouter(e *echo.Echo) {
	g := e.Group("/message")
	g.POST("", controller.SendMessage, middleware.JWTWithConfig(util.Conf))
}
