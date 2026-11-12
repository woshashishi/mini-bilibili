package like

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"mini-bilibili/internal/video"
)

var (
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrInvalidVideoID = errors.New("invalid video id")
	ErrVideoNotFound  = errors.New("video not found")
)

type Service struct {
	repo      *Repository
	videoRepo *video.Repository
}

func NewService(
	repo *Repository,
	videoRepo *video.Repository,
) *Service {
	// TODO：保存两个依赖，返回 Service
	return &Service{
		repo:      repo,
		videoRepo: videoRepo,
	}
}

func (s *Service) Add(
	ctx context.Context,
	userID uint64,
	videoID uint64,
) error {
	// TODO 1：userID 为零时，返回 ErrInvalidUserID
	if userID == 0 {
		return ErrInvalidUserID
	}
	// TODO 2：videoID 为零时，返回 ErrInvalidVideoID
	if videoID == 0 {
		return ErrInvalidVideoID
	}
	// TODO 3：调用 videoRepo.FindByID
	// 不需要视频对象，可以用 _ 接收第一个返回值
	_, err := s.videoRepo.FindByID(ctx, videoID)
	// TODO 4：处理查询错误
	// 匹配 gorm.ErrRecordNotFound → 返回 ErrVideoNotFound
	if err != nil {
		ok := errors.Is(err, gorm.ErrRecordNotFound)
		if ok {
			return ErrVideoNotFound
		}
		return err
	}
	// 其他非空错误 → 原样返回
	// 使用 errors.Is 判断
	// TODO 5：调用点赞 Repository.Add
	// 传入 ctx、userID、videoID
	// 返回它的执行结果
	err = s.repo.Add(ctx, userID, videoID)
	if err != nil {
		return err
	}
	return nil
}
