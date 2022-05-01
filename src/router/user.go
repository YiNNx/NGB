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

	g.GET("/:uid", controller.GetUserProfile)
	g.PUT("/follow/:uid", controller.FollowUser, middleware.JWTWithConfig(util.Conf))

	g.GET("/account", controller.GetUserAccount, middleware.JWTWithConfig(util.Conf))
	g.PUT("/account", controller.ChangeUserInfo, middleware.JWTWithConfig(util.Conf))
	g.PUT("/password", controller.ChangeUserPwd, middleware.JWTWithConfig(util.Conf))

	g.GET("/all", controller.GetAllUsers, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.GET("/admin", controller.GetAdmins, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
	g.DELETE("/:id", controller.DeleteUser, middleware.JWTWithConfig(util.Conf), myware.VerifySuperAdmin)
}
