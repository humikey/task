package router

import (
	"github.com/gin-gonic/gin"
	"my-blog/controller"
	"my-blog/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.GET("/users/:id/posts", controller.GetUserPosts)
	r.GET("/posts/most_commented", controller.GetMostCommentedPost)
	api := r.Group("/api")
	{
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
		api.GET("/posts/:id", controller.GetPostDetail) // ✅ 公共接口，任何人都能看

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/posts", controller.CreatePost)                 // 需要登录才可发帖
			auth.POST("/posts/:id/comments", controller.CreateComment) // 需要登录才可评论
			auth.PUT("/posts/:id", controller.UpdatePost)              // 需要登录才可修改
			auth.DELETE("/posts/:id", controller.DeletePost)           // 需要登录才可删除
		}
	}
	return r
}
