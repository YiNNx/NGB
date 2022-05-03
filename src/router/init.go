package router

import (
	"github.com/labstack/echo/v4"
)

func InitRouters(e *echo.Echo) {
	initUserRouter(e)
	initPostRouter(e)
	initBoardRouter(e)
	initApplyRouter(e)
	initNotificationRouter(e)
	initMessageRouter(e)
}
