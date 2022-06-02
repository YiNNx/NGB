package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initApplyRouter(e *echo.Echo) {
	g := e.Group("/apply")

	g.POST("/admin", controller.SetAdminApply, myware.HandleSession)
	g.POST("/board", controller.SetBoardApply, myware.HandleSession)

	g.GET("/admin", controller.GetAdminApply, myware.HandleSession, myware.VerifySuperAdmin)
	g.GET("/board", controller.GetBoardApply, myware.HandleSession, myware.VerifySuperAdmin)
	g.POST("/:apid", controller.PassApply, myware.HandleSession, myware.VerifySuperAdmin)
}
