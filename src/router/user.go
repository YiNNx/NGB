package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initUserRouter(e *echo.Echo) {
	g := e.Group("/user")

	g.POST("", controller.SignUp)
	g.GET("/token", controller.LogIn)
	g.POST("/logout", controller.LogOut)
	g.GET("/:uid", controller.GetUserProfile)

	g.PUT("/follow/:uid", controller.FollowUser, myware.HandleSession)
	g.GET("/account", controller.GetUserAccount, myware.HandleSession)
	g.PUT("/account", controller.ChangeUserInfo, myware.HandleSession)
	g.PUT("/password", controller.ChangeUserPwd, myware.HandleSession)

	g.GET("/all", controller.GetAllUsers, myware.HandleSession, myware.VerifySuperAdmin)
	g.GET("/admin", controller.GetAdmins, myware.HandleSession, myware.VerifySuperAdmin)
	g.DELETE("/:id", controller.DeleteUser, myware.HandleSession, myware.VerifySuperAdmin)
}
