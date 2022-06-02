package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initNotificationRouter(e *echo.Echo) {
	g := e.Group("/notification")

	g.GET("", controller.GetNotification, myware.HandleSession)
	g.GET("/new", controller.GetNewNotification, myware.HandleSession)
}
