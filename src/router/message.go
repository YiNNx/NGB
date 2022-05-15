package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	"ngb/util"
)

func initMessageRouter(e *echo.Echo) {
	e.GET("/chat", controller.Chat, middleware.JWTWithConfig(util.Conf))
	e.POST("/message", controller.SendMessage, middleware.JWTWithConfig(util.Conf))
}
