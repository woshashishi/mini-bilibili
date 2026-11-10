package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mini-bilibili/internal/middleware"
	"strconv"
	"time"
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
type DetailResponse struct {
	ID          uint64    `json:"id"`
	AuthorID    uint64    `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type FeedResponse struct {
	Items      []DetailResponse `json:"items"`
	HasMore    bool             `json:"has_more"`
	NextCursor uint64           `json:"next_cursor"`
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

func (h *Handler) GetByID(c *gin.Context) {
	// TODO 1：通过 c.Param("id") 读取路径参数
	str := c.Param("id")
	// TODO 2：用 strconv.ParseUint 转成 uint64
	// 参数：路径字符串、10、64
	// 转换失败：
	// HTTP 400、code 40004、message "视频ID格式错误"
	// 然后 return
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40004,
			"message": "视频ID格式错误",
		})
		return
	}
	// TODO 3：调用 h.service.GetByID
	// 传入 c.Request.Context() 和转换后的 ID
	video, err := h.service.GetByID(c.Request.Context(), num)
	// TODO 4：使用 errors.Is 处理错误
	//
	// ErrInvalidVideoID：
	// HTTP 400、code 40004、message "视频ID必须大于零"
	//
	// ErrVideoNotFound：
	// HTTP 404、code 40401、message "视频不存在"
	//
	// 其他错误：
	// HTTP 500、code 50001、message "服务器内部错误"
	//
	// 每个错误响应后 return
	if err != nil {
		if ok := errors.Is(err, ErrInvalidVideoID); ok {
			c.JSON(400, gin.H{
				"code":    40004,
				"message": "视频ID必须大于零",
			})
			return
		}
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
	// TODO 5：将 Video 中的字段填入 DetailResponse
	// 返回 HTTP 200
	// code：0
	// message："获取成功"
	// data：详情响应对象
	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": DetailResponse{
			ID:          video.ID,
			AuthorID:    video.AuthorID,
			Title:       video.Title,
			VideoURL:    video.VideoURL,
			CreatedAt:   video.CreatedAt,
			Description: video.Description,
		},
	})

}
func (h *Handler) Feed(c *gin.Context) {
	// TODO 1：读取 cursor，默认字符串 "0"
	// 用 strconv.ParseUint(rawCursor, 10, 64) 转换
	// 失败：HTTP 400、code 40005
	// message："游标格式错误"
	// 然后 return
	cursorStr := c.DefaultQuery("cursor", "0")
	var cursor uint64 = 0
	var err error
	cursor, err = strconv.ParseUint(cursorStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40005,
			"message": "游标格式错误",
		})
		return
	}
	// TODO 2：读取 limit，默认字符串 "10"
	// 用 strconv.Atoi 转换为 int
	// 失败：HTTP 400、code 40005
	// message："分页数量格式错误"
	// 然后 return
	limitStr := c.DefaultQuery("limit", "10")
	var limit int = 10
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    40005,
			"message": "分页数量格式错误",
		})
		return
	}

	// TODO 3：调用 h.service.Feed
	// 传入请求 Context、cursor、limit
	feedResult, err := h.service.Feed(c.Request.Context(), cursor, limit)
	// TODO 4：处理错误
	// ErrInvalidLimit：
	// HTTP 400、code 40005
	// message："分页数量应为0至50，0表示默认10条"
	//
	// 其他错误：
	// HTTP 500、code 50001
	// message："服务器内部错误"
	//
	// 每个错误响应后 return
	if err != nil {
		if errors.Is(err, ErrInvalidLimit) {
			c.JSON(400, gin.H{
				"code":    40005,
				"message": "分页数量应为0至50，0表示默认10条",
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
	// TODO 5：创建非 nil 的空 []DetailResponse
	// 遍历 result.Items，将每条视频转换为 DetailResponse
	// 包含 ID、AuthorID、Title、Description、VideoURL、CreatedAt
	detailResponse := make([]DetailResponse, 0, len(feedResult.Items))
	for i := 0; i < len(feedResult.Items); i++ {
		item := DetailResponse{
			ID:          feedResult.Items[i].ID,
			AuthorID:    feedResult.Items[i].AuthorID,
			Title:       feedResult.Items[i].Title,
			Description: feedResult.Items[i].Description,
			VideoURL:    feedResult.Items[i].VideoURL,
			CreatedAt:   feedResult.Items[i].CreatedAt,
		}
		detailResponse = append(detailResponse, item)
	}
	// TODO 6：返回 HTTP 200
	// code：0
	// message："获取成功"
	// data：FeedResponse
	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": &FeedResponse{
			Items:      detailResponse,
			HasMore:    feedResult.HasMore,
			NextCursor: feedResult.NextCursor,
		},
	})
}
