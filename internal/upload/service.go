package upload

import (
	"io"

	"mini-bilibili/internal/storage"
)

type Service struct {
	store *storage.LocalStorage
}

func NewService(store *storage.LocalStorage) *Service {
	// TODO：保存存储对象，返回 Service
	return &Service{
		store: store,
	}
}

func (s *Service) SaveVideo(
	filename string,
	size int64,
	src io.ReadSeeker,
) (string, error) {
	// TODO 1：调用 ValidateMP4
	// 传入 filename、size、src
	// 校验失败：返回空字符串和原始错误
	err := ValidateMP4(filename, size, src)
	if err != nil {
		return "", err
	}

	// TODO 2：调用 s.store.SaveMP4(src)
	// 保存失败：返回空字符串和原始错误
	str, err := s.store.SaveMP4(src)
	// TODO 3：成功返回生成的文件名和 nil
	if err != nil {
		return "", err
	}
	return str, err
}
