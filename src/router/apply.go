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

	g.GET("/", controller.GetApply, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.POST("/:apid", controller.PassApply, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)

	g.POST("/admin", controller.SetAdminApply, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)
	g.POST("/board", controller.SetBoardApply, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)
}
