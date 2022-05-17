package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	"ngb/util"
)

func initEmailRouter(e *echo.Echo) {
	e.POST("/admin/email", controller.SendEmail, middleware.JWTWithConfig(util.Conf))
}
