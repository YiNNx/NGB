package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	myware "ngb/middleware"
	"ngb/util"
)

func initApplyRouter(e *echo.Echo) {
	g := e.Group("/apply")

	g.POST("/admin", controller.SetAdminApply, middleware.JWTWithConfig(util.Conf))
	g.POST("/board", controller.SetBoardApply, middleware.JWTWithConfig(util.Conf))

	g.GET("/admin", controller.GetAdminApply, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.GET("/board", controller.GetBoardApply, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.POST("/:apid", controller.PassApply, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
}
