package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/controller"
	myware "ngb/middleware"
	"ngb/util"
)

func initPostRouter(e *echo.Echo) {
	e.POST("/board/:bid/post", controller.NewPost, middleware.JWTWithConfig(util.Conf))

	g := e.Group("/post")

	g.GET("", controller.GetPostsByTag)
	g.GET("/all", controller.GetAllPosts)
	g.GET("/:pid", controller.GetPost)

	g.PUT("/:pid/collection", controller.CollectPost, middleware.JWTWithConfig(util.Conf))
	g.PUT("/:pid/like", controller.LikePost, middleware.JWTWithConfig(util.Conf))
	g.POST("/:pid/comment", controller.CommentPost, middleware.JWTWithConfig(util.Conf))
	g.POST("/:pid/comment/:cid/subcomment", controller.SubCommentPost, middleware.JWTWithConfig(util.Conf))

	g.DELETE("/:pid", controller.DeletePost, middleware.JWTWithConfig(util.Conf), myware.VerifyAdmin)

}
