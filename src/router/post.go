package router

import (
	"github.com/labstack/echo/v4"
	"ngb/controller"
	myware "ngb/middleware"
)

func initPostRouter(e *echo.Echo) {
	e.POST("/board/:bid/post", controller.NewPost, myware.HandleSession)

	g := e.Group("/post")

	g.GET("", controller.GetPostsByTag)
	g.GET("/search", controller.SearchPost)
	g.GET("/all", controller.GetAllPosts)
	g.GET("/:pid", controller.GetPost)

	g.PUT("/:pid/collection", controller.CollectPost, myware.HandleSession)
	g.PUT("/:pid/like", controller.LikePost, myware.HandleSession)
	g.POST("/:pid/comment", controller.CommentPost, myware.HandleSession)
	g.POST("/:pid/comment/:cid/subcomment", controller.SubCommentPost, myware.HandleSession)

	g.DELETE("/:pid", controller.DeletePost, myware.HandleSession, myware.VerifyAdmin)

}
