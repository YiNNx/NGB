package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
)

func initBoardRouter(e *echo.Echo) {
	g := e.Group("/board")
	g.GET("/all", controller.GetAllBoards)
	g.GET("/:bid", controller.GetBoard)
}
