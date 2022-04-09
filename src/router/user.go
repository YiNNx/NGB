package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	"ngb/util"
)

func initUserRouter(e *echo.Echo) {
	g := e.Group("/user")
	g.POST("", controller.SignUP)
	g.GET("/token", controller.LogIn)

	g.GET("/:uid", controller.GetUser, middleware.JWTWithConfig(util.Conf), util.VerifyUser)
	g.PUT("/:uid", controller.ChangeInfo, middleware.JWTWithConfig(util.Conf), util.VerifyUser)
	g.PUT("/:uid/pwd", controller.ChangePwd, middleware.JWTWithConfig(util.Conf), util.VerifyUser)
	g.PUT("/:uid/follow", controller.Follow, middleware.JWTWithConfig(util.Conf), util.VerifyUser)

	g.GET("/all", controller.GetAllUsers, middleware.JWTWithConfig(util.Conf), util.VerifyAdmin)
	g.DELETE("/:id", controller.DeleteUser, middleware.JWTWithConfig(util.Conf), util.VerifyAdmin)
}
