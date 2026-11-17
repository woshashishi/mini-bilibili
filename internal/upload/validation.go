package upload

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const MaxVideoSize int64 = 100 * 1024 * 1024

var (
	ErrInvalidSize = errors.New("invalid video size")
	ErrInvalidType = errors.New("invalid video type")
)

func ValidateMP4(
	filename string,
	size int64,
	src io.ReadSeeker,
) error {
	// TODO 1：size <= 0 或超过 MaxVideoSize
	// 返回 ErrInvalidSize
	if size <= 0 || size > MaxVideoSize {
		return ErrInvalidSize
	}
	// TODO 2：用 filepath.Ext 获取扩展名
	// 用 strings.EqualFold 与 ".mp4" 比较
	// 不符合时返回 ErrInvalidType
	str := filepath.Ext(filename)
	ok := strings.EqualFold(str, ".mp4")
	if !ok {
		return ErrInvalidType
	}
	// TODO 3：使用 src.Seek(0, io.SeekStart)
	// 将读取位置移到文件开头
	// 失败时返回原始错误
	_, err := src.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// TODO 4：创建 512 字节缓冲区
	// 使用 io.ReadFull，接收 n 和 err
	//
	var buf = make([]byte, 512)
	// io.EOF、io.ErrUnexpectedEOF：
	// 表示文件不足512字节，继续检测已读取内容
	//
	n, err := io.ReadFull(src, buf)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return err
	}
	// 其他非空错误：
	// 返回原始错误
	str = http.DetectContentType(buf[:n])
	if str != "video/mp4" {
		return ErrInvalidType
	}
	// TODO 5：调用 http.DetectContentType(buf[:n])
	// 结果不是 "video/mp4" 时返回 ErrInvalidType

	// TODO 6：再次 Seek 到文件开头
	// 检查 Seek 错误，然后返回 nil
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
