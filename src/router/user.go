package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	myware "ngb/middleware"
	"ngb/util"
)

func initUserRouter(e *echo.Echo) {
	g := e.Group("/user")

	g.POST("", controller.SignUP)
	g.GET("/token", controller.LogIn)

	g.GET("/:uid", controller.GetUserProfile, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.PUT("/user/account", controller.GetUserAccount, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.PUT("/user/account", controller.ChangeUserInfo, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.PUT("/user/pwd", controller.ChangeUserPwd, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.PUT("/follow/:uid", controller.FollowUser, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)

	g.GET("/all", controller.GetAllUsers, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)
	g.DELETE("/:id", controller.DeleteUser, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)
}
