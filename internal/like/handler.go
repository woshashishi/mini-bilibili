package like

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/middleware"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	// TODO：保存 Service
	return &Handler{
		service: service,
	}
}

func (h *Handler) Add(c *gin.Context) {
	// TODO 1：读取 middleware.UserIDKey
	// 检查存在、安全断言为 uint64，并检查大于零
	//
	val, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// 不符合要求：
	// HTTP 500、code 50001、message "服务器内部错误"
	// 然后 return
	userId, ok := val.(uint64)
	if !ok || userId <= 0 {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// TODO 2：读取 c.Param("id")
	// 使用 strconv.ParseUint 转换
	//
	// 转换失败：
	// HTTP 400、code 40004、message "视频ID格式错误"
	// 然后 return
	str := c.Param("id")
	videoId, err := strconv.ParseUint(str, 10, 64)
	// TODO 3：调用 h.service.Add
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40004,
			"message": "视频ID格式错误",
		})
		return
	}
	// 传入 c.Request.Context()、用户 ID、视频 ID
	err = h.service.Add(c.Request.Context(), userId, videoId)
	// TODO 4：用 errors.Is 判断业务错误
	//
	// ErrInvalidVideoID：
	// HTTP 400、code 40004、message "视频ID必须大于零"
	//
	if err != nil {
		if ok := errors.Is(err, ErrInvalidVideoID); ok {
			c.JSON(400, gin.H{
				"code":    40004,
				"message": "视频ID必须大于零",
			})
			return
		}
		// ErrVideoNotFound：
		// HTTP 404、code 40401、message "视频不存在"
		//
		if ok := errors.Is(err, ErrVideoNotFound); ok {
			c.JSON(404, gin.H{
				"code":    40401,
				"message": "视频不存在",
			})
			return
		}
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// 其他错误，包括 ErrInvalidUserID：
	// HTTP 500、code 50001、message "服务器内部错误"
	//
	// 每个错误响应后 return
	// TODO 5：成功返回 HTTP 200
	// code：0
	// message："点赞成功"
	c.JSON(200, gin.H{
		"code":    0,
		"message": "点赞成功",
	})
}

func (h *Handler) Remove(c *gin.Context) {
	// TODO 1：读取 middleware.UserIDKey
	// 检查存在、安全断言为 uint64、检查非零
	// 失败：HTTP 500、code 50001
	// message："服务器内部错误"
	// 然后 return
	val, exit := c.Get(middleware.UserIDKey)
	if !exit {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	userId, ok := val.(uint64)
	if !ok || userId == 0 {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// TODO 2：解析 c.Param("id")
	// 使用 strconv.ParseUint，十进制、64位
	// 转换失败：HTTP 400、code 40004
	// message："视频ID格式错误"
	// 然后 return
	str := c.Param("id")
	videoId, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40004,
			"message": "视频ID格式错误",
		})
		return
	}
	// TODO 3：调用 h.service.Remove
	// 传入请求 Context、用户 ID、视频 ID
	// 必须接收返回的 error
	err = h.service.Remove(c.Request.Context(), userId, videoId)
	// TODO 4：处理错误
	// ErrInvalidVideoID：
	// HTTP 400、code 40004、message "视频ID必须大于零"
	//
	if err != nil {
		if errors.Is(err, ErrInvalidVideoID) {
			c.JSON(400, gin.H{
				"code":    40004,
				"message": "视频ID必须大于零",
			})
			return
		}
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// 其他错误：
	// HTTP 500、code 50001、message "服务器内部错误"
	//
	// 每个错误响应后 return

	// TODO 5：成功返回 HTTP 200
	// code：0
	// message："取消点赞成功"
	c.JSON(200, gin.H{
		"code":    0,
		"message": "取消点赞成功",
	})
	return
}
