package server

import (
	"github.com/gin-gonic/gin"

	"mini-bilibili/internal/user"
)

func NewRouter(userHandler *user.Handler) *gin.Engine {
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
	return router
}
