package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initBoardRouter(e *echo.Echo) {
	g := e.Group("/board")
	g.GET("/all", controller.GetAllBoards)
	g.GET("/:bid", controller.GetBoard)

	g.POST("/:bid", controller.UpdateBoard, myware.HandleSession, myware.VerifyAdmin)

	g.POST("", controller.SetBoard, myware.HandleSession, myware.VerifySuperAdmin)
	g.DELETE("/:bid", controller.DeleteBoard, myware.HandleSession, myware.VerifySuperAdmin)
}
