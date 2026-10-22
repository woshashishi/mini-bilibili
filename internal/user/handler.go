package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	// TODO：创建并返回 Handler
	return &Handler{
		service: service,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// TODO 1：使用 ShouldBindJSON 解析请求
	// 解析失败时返回：
	// HTTP 400
	// {
	//   "code": 40001,
	//   "message": "请求参数错误"
	// }
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40001,
			"message": "请求参数错误",
		})
		return
	}
	// TODO 2：调用 h.service.Register
	u, err := h.service.Register(c.Request.Context(), req.Username, req.Password)
	// context 使用 c.Request.Context()
	// TODO 3：处理 Service 返回的错误
	// ErrInvalidUsername、ErrInvalidPassword：
	// HTTP 400，code 为 40002
	// ErrUsernameExists：
	// HTTP 409，code 为 40901
	// 其他错误：
	// HTTP 500，code 为 50001

	// 判断错误使用 errors.Is
	if err != nil {
		if errors.Is(err, ErrInvalidUsername) || errors.Is(err, ErrInvalidPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    40002,
				"message": "用户名或密码格式错误",
			})
			return
		}
		if errors.Is(err, ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    40901,
				"message": "用户名已存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50001,
			"message": "服务器内部错误", // 补充 message
		})
		return
	}
	// 必须 return，防止往下走触发 u.ID 空指针
	// TODO 4：成功时返回 HTTP 201
	// code 为 0
	// message 为 "注册成功"
	// data 使用 RegisterResponse，只返回 ID 和 Username
	rep := &RegisterResponse{
		ID:       uint64(u.ID),
		Username: u.Username,
	}
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "注册成功",
		"data":    rep,
	})
}
