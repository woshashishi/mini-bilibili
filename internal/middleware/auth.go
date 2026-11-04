package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"mini-bilibili/internal/auth"
)

const UserIDKey = "user_id"

func RequireAuth(manager *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 1：通过 c.GetHeader 读取 Authorization
		header := c.GetHeader("Authorization")
		// TODO 2：使用 strings.Fields 拆分
		// 必须恰好有两部分
		// 第一部分用 strings.EqualFold 判断是否为 "Bearer"
		//
		st := strings.Fields(header)
		if len(st) != 2 || !strings.EqualFold(st[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    40102,
				"message": "请提供有效的登录凭证",
			})
			return
		}
		// 格式不正确：
		// 使用 c.AbortWithStatusJSON
		// HTTP：http.StatusUnauthorized
		// code：40102
		// message："请提供有效的登录凭证"
		// 然后 return

		// TODO 3：将第二部分交给 manager.Parse
		//
		cl, err := manager.Parse(st[1])
		// 验证失败：
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    40102,
				"message": "登录凭证无效或已过期",
			})
			return
		}
		// 使用 c.AbortWithStatusJSON
		// HTTP：http.StatusUnauthorized
		// code：40102
		// message:"登录凭证无效或已过期"
		// 然后 return
		c.Set(UserIDKey, cl.UserID)
		// TODO 4：将验证后的 Claims.UserID 保存到 Context
		// 使用 c.Set，键为 UserIDKey
		c.Next()
		// TODO 5：使用 c.Next() 执行后续 Handler
	}
}
