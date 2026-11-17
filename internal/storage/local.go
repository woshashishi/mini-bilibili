package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	root string
}

func NewLocalStorage(root string) (*LocalStorage, error) {
	// TODO 1：root 为空时，返回配置错误
	if root == "" {
		return nil, errors.New("配置错误")
	}
	// TODO 2：使用 os.MkdirAll(root, 0755)
	// 确保存储目录存在
	// 失败时返回错误
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return nil, err
	}
	// TODO 3：保存 root，返回 LocalStorage
	return &LocalStorage{
		root: root,
	}, nil
}

// SaveMP4 接收已经通过上层校验的文件流。
// 返回服务器生成的文件名，不返回本机绝对路径。
func (s *LocalStorage) SaveMP4(src io.Reader) (string, error) {
	// TODO 1：创建长度为 16 的随机字节切片
	// 用 rand.Read 填充，必须处理错误
	var b = make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// TODO 2：hex.EncodeToString 转成字符串
	// 拼接 ".mp4"，得到 filename
	str := hex.EncodeToString(b)
	// TODO 3：filepath.Join(s.root, filename)
	// 得到磁盘目标路径
	filename := str + ".mp4"
	fullpath := filepath.Join(s.root, filename)
	// TODO 4：使用 os.OpenFile 创建文件
	// 标志：os.O_WRONLY | os.O_CREATE | os.O_EXCL
	// 权限：0600
	// 失败时返回错误
	file, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return "", err
	}
	// TODO 5：使用 io.Copy，将 src 写入文件
	// 无论复制是否成功，都必须关闭文件
	//
	_, err = io.Copy(file, src)
	// 复制或关闭失败：
	// 尝试删除不完整文件，并返回错误
	//
	if err != nil {
		file.Close()
		os.Remove(fullpath)
		return "", err
	}
	// 注意：删除前先关闭文件，尤其在 Windows 上
	err = file.Close()
	if err != nil {
		os.Remove(fullpath)
		return "", err
	}
	// TODO 6：成功返回 filename 和 nil
	return filename, nil
}
