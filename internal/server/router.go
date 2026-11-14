package server

import (
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/auth"
	"mini-bilibili/internal/like"
	"mini-bilibili/internal/middleware"
	"mini-bilibili/internal/video"

	"mini-bilibili/internal/user"
)

func NewRouter(userHandler *user.Handler, tokenManager *auth.TokenManager, videoHandler *video.Handler, likeHandler *like.Handler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// TODO：建立路由分组
	// 最终地址应该是：
	// POST /api/v1/auth/register
	api := router.Group("api/v1")
	auth := api.Group("auth")
	auth.POST("register", userHandler.Register)
	auth.POST("login", userHandler.Login)
	// TODO：注册 GET /me
	// 处理顺序：
	// middleware.RequireAuth(tokenManager)
	// userHandler.Me
	api.GET("/me", middleware.RequireAuth(tokenManager), userHandler.Me)
	// TODO：注册 POST /videos
	// 第一个处理函数：middleware.RequireAuth(tokenManager)
	// 第二个处理函数：videoHandler.Publish
	api.POST("/videos", middleware.RequireAuth(tokenManager), videoHandler.Publish)

	// TODO：注册 GET /videos/:id
	// 直接绑定 videoHandler.GetByID
	// 第一版视频详情允许公开查看，不挂鉴权中间件
	api.GET("/videos/:id", videoHandler.GetByID)
	// TODO：注册 GET /feed，绑定 videoHandler.Feed
	// 第一版 Feed 允许未登录访问，不挂鉴权中间件
	api.GET("/feed", videoHandler.Feed)
	// TODO：使用 api.PUT 注册 "/videos/:id/like"
	// 先执行 middleware.RequireAuth(tokenManager)
	// 再执行 likeHandler.Add
	api.PUT("/videos/:id/like", middleware.RequireAuth(tokenManager), likeHandler.Add)

	// TODO：使用 DELETE 注册相同路径
	// 先鉴权，再调用 likeHandler.Remove
	api.DELETE("/videos/:id/like", middleware.RequireAuth(tokenManager), likeHandler.Remove)
	// TODO：注册 GET /videos/:id/likes/count
	// 绑定 likeHandler.Count，不挂鉴权中间件
	api.GET("/videos/:id/likes/count", likeHandler.Count)
	// TODO：注册 GET /videos/:id/like
	// 先执行 middleware.RequireAuth(tokenManager)
	// 再执行 likeHandler.Status
	api.GET("/videos/:id/like", middleware.RequireAuth(tokenManager), likeHandler.Status)
	return router
}
