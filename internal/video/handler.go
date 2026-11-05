package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	// TODO：保存 Service，返回 Handler
	return &Handler{
		service: service,
	}
}

type PublishRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url" binding:"required"`
}

type PublishResponse struct {
	ID uint64 `json:"id"`
}

func (h *Handler) Publish(c *gin.Context) {
	// TODO 1：读取 middleware.UserIDKey
	val, exits := c.Get(middleware.UserIDKey)
	if !exits {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	userID, ok := val.(uint64)
	if !ok || userID == 0 {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// 使用 c.Get，并安全断言为 uint64
	// 不存在、类型错误或 ID 为零：
	// HTTP 500、code 50001、message "服务器内部错误"
	// 然后 return

	var req PublishRequest

	// TODO 2：ShouldBindJSON 解析请求
	// 失败：
	// HTTP 400、code 40001、message "请求参数错误"
	// 然后 return
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40001,
			"message": "请求参数错误",
		})
		return
	}
	// TODO 3：调用 h.service.Publish
	// 依次传入：
	// c.Request.Context()
	video, err := h.service.Publish(c.Request.Context(), userID, req.Title, req.Description, req.VideoURL)
	// 鉴权得到的用户 ID
	// req.Title
	// req.Description
	// req.VideoURL
	if err != nil {
		if errors.Is(err, ErrInvalidTitle) {
			c.JSON(400, gin.H{
				"code":    40003,
				"message": "标题必须为1至100个字符",
			})
			return
		} else if errors.Is(err, ErrInvalidDescription) {
			c.JSON(400, gin.H{
				"code":    40003,
				"message": "简介不能超过1000个字符",
			})
			return
		} else if errors.Is(err, ErrInvalidVideoURL) {
			c.JSON(400, gin.H{
				"code":    40003,
				"message": "视频地址格式错误或过长",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"code":    50001,
				"message": "服务器内部错误",
			})
			return
		}
	}

	c.JSON(201, gin.H{
		"code":    0,
		"message": "发布成功",
		"data": &PublishResponse{
			ID: video.ID,
		},
	})
	// TODO 4：处理业务错误，使用 errors.Is
	//
	// ErrInvalidTitle：
	// HTTP 400、code 40003、message "标题必须为1至100个字符"
	//
	// ErrInvalidDescription：
	// HTTP 400、code 40003、message "简介不能超过1000个字符"
	//
	// ErrInvalidVideoURL：
	// HTTP 400、code 40003、message "视频地址格式错误或过长"
	//
	// 其他错误，包括 ErrInvalidAuthor：
	// HTTP 500、code 50001、message "服务器内部错误"
	//
	// 每个错误响应后 return

	// TODO 5：成功返回 HTTP 201
	// code：0
	// message："发布成功"
	// data：PublishResponse，ID 为创建成功的视频 ID

}
