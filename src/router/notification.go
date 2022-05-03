package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	"ngb/util"
)

func initNotificationRouter(e *echo.Echo) {
	g := e.Group("/notification")
	g.GET("", controller.GetNotification, middleware.JWTWithConfig(util.Conf))
	g.GET("/new", controller.GetNewNotification, middleware.JWTWithConfig(util.Conf))
}
