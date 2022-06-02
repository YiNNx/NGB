package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initEmailRouter(e *echo.Echo) {
	e.POST("/admin/email", controller.SendEmail, myware.HandleSession, myware.VerifySuperAdmin)
}
