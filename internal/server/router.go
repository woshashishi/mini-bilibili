package server

import (
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/auth"
	"mini-bilibili/internal/middleware"

	"mini-bilibili/internal/user"
)

func NewRouter(userHandler *user.Handler, tokenManager *auth.TokenManager) *gin.Engine {
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

	return router
}
