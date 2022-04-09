package router

import "github.com/labstack/echo/v4"

func initPostRouter(e *echo.Echo) {
	g := e.Group("/post")
	g.POST("", controller.Post)
	g.GET("/all", controller.GetAllPosts)
	g.GET("/:pid", controller.GetPost)
	g.PUT("/post/:pid/collection", controller.CollectPost())
	g.PUT("/post/:pid/like", controller.LikePost())
	g.POST("/post/:pid/comment", controller.CommentPost())
}
