package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	myware "ngb/middleware"
	"ngb/util"
)

func initPostRouter(e *echo.Echo) {
	g := e.Group("/post")
	g.POST("", controller.NewPost, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.GET("/all", controller.GetAllPosts)
	g.GET("/:pid", controller.GetPost)
	g.PUT("/post/:pid/collection", controller.CollectPost, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.PUT("/post/:pid/like", controller.LikePost, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
	g.POST("/post/:pid/comment", controller.CommentPost, middleware.JWTWithConfig(util.Conf), myware.VerifyUser)
}
