package upload

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"mini-bilibili/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	// TODO：保存 service，返回 Handler
	return &Handler{
		service: service,
	}
}

type UploadResponse struct {
	Filename string `json:"filename"`
}

func (h *Handler) UploadVideo(c *gin.Context) {
	// TODO 1：读取 middleware.UserIDKey
	// 检查存在、类型为 uint64、值大于零
	// 不符合要求：返回 500、code 50001，然后 return
	val, exit := c.Get(middleware.UserIDKey)
	if !exit {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	userId, ok := val.(uint64)
	if !ok || userId <= 0 {
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// 在解析文件之前，限制整个请求体。
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		MaxVideoSize+(1<<20),
	)

	err := c.Request.ParseMultipartForm(8 << 20)

	// 先注册清理，解析后续出错也能清理已有临时文件。
	defer func() {
		if c.Request.MultipartForm != nil {
			_ = c.Request.MultipartForm.RemoveAll()
		}
	}()

	// TODO 2：处理解析错误
	// 使用 errors.As 判断 *http.MaxBytesError：
	//   HTTP 413、code 41301、message "上传请求过大"
	//
	var sizeErr *http.MaxBytesError
	if err != nil {
		ok = errors.As(err, &sizeErr)
		if ok {
			c.JSON(413, gin.H{
				"code":    41301,
				"message": "上传请求过大",
			})
			return
		}
		c.JSON(400, gin.H{
			"code":    40006,
			"message": "上传表单格式错误",
		})
		return

	}

	// 其他解析错误：
	//   HTTP 400、code 40006、message "上传表单格式错误"
	//
	// 返回响应后立即 return

	file, header, err := c.Request.FormFile("file")

	// TODO 3：处理获取文件的错误
	// http.ErrMissingFile：
	//   HTTP 400、code 40006、message "请选择视频文件"
	//
	if err != nil {
		ok := errors.Is(err, http.ErrMissingFile)
		if ok {
			c.JSON(400, gin.H{
				"code":    40006,
				"message": "请选择视频文件",
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
	//   HTTP 500、code 50001、message "服务器内部错误"
	//
	// 返回响应后立即 return

	// 必须在确认获取文件成功后注册。
	defer file.Close()

	// TODO 4：空文件返回 400、code 40006
	// message "视频文件不能为空"
	// 使用 header.Size 判断，然后 return
	if header.Size == 0 {
		c.JSON(400, gin.H{
			"code":    40006,
			"message": "视频文件不能为空",
		})
		return
	}
	// TODO 5：调用 h.service.SaveVideo
	// 依次传入 header.Filename、header.Size、file
	// 接收服务器生成的文件名和错误
	str, err := h.service.SaveVideo(c.Request.Context(), header.Filename, header.Size, file)
	// TODO 6：处理 Service 错误
	// ErrInvalidSize：
	//   HTTP 413、code 41301、message "视频不能超过100MiB"
	//
	if err != nil {
		if errors.Is(err, ErrInvalidSize) {
			c.JSON(413, gin.H{
				"code":    41301,
				"message": "视频不能超过100MiB",
			})
			return
		}
		if errors.Is(err, ErrInvalidType) {
			c.JSON(400, gin.H{
				"code":    40006,
				"message": "请上传支持的MP4文件",
			})
			return
		}
		c.JSON(500, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
		return
	}
	// ErrInvalidType：
	//   HTTP 400、code 40006、message "请上传支持的MP4文件"
	//
	// 其他错误：
	//   HTTP 500、code 50001、message "服务器内部错误"
	//
	// 返回响应后立即 return

	// TODO 7：成功返回 HTTP 201
	// code：0
	// message："上传成功"
	// data：UploadResponse，填入生成的文件名
	c.JSON(201, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": &UploadResponse{
			Filename: str,
		},
	})
	return
}
