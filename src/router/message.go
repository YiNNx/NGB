package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initMessageRouter(e *echo.Echo) {
	e.GET("/chat", controller.Chat, myware.HandleSession, myware.WsUpgrade)
	e.POST("/message", controller.SendMessage, myware.HandleSession)
}
