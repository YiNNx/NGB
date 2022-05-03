package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	myware "ngb/middleware"
	"ngb/util"
)

func initBoardRouter(e *echo.Echo) {
	g := e.Group("/board")
	g.GET("/all", controller.GetAllBoards)
	g.GET("/:bid", controller.GetBoard)

	g.POST("/:bid", controller.UpdateBoard, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)

	g.POST("", controller.SetBoard, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.DELETE("/:bid", controller.DeleteBoard, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
}
