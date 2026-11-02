package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/auth"
	"net/http"
)

type Handler struct {
	service      *Service
	tokenManager *auth.TokenManager
}

func NewHandler(service *Service, tokenManager *auth.TokenManager) *Handler {
	// TODO：创建并返回 Handler
	return &Handler{
		service:      service,
		tokenManager: tokenManager,
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
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
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
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// TODO 1：ShouldBindJSON 解析请求
	// 失败：HTTP 400
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40001,
			"message": "请求参数错误",
		})
		return
	}

	// code：40001
	// message:"请求参数错误"
	// 返回响应后立即 return

	// TODO 2：调用 h.service.Login
	// 传入请求的 Context、用户名和密码
	u, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	// TODO 3：处理登录错误
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(401, gin.H{
				"code":    40101,
				"message": "用户名或密码错误",
			})
			return
		}
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}

	// ErrInvalidCredentials：
	//   HTTP 401，code 40101
	//   message："用户名或密码错误"
	//
	// 其他错误：
	//   HTTP 500，code 50001
	//   message："服务器内部错误"
	//
	// 每个错误响应后立即 return

	// TODO 4：使用登录成功的用户 ID 调用 Generate
	// 签发失败：返回 HTTP 500，然后 return

	st, err := h.tokenManager.Generate(uint64(u.ID))
	// TODO 5：返回 HTTP 200
	if err != nil {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	lp := &LoginResponse{
		AccessToken: st,
		TokenType:   "Bearer",
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    lp,
	})
	// code：0
	// message："登录成功"
	// data：LoginResponse
	//   AccessToken：签发的 Token
	//   TokenType："Bearer"
}
