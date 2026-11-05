package video

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidAuthor      = errors.New("invalid author")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidVideoURL    = errors.New("invalid video url")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	// TODO：保存 Repository，返回 Service
	return &Service{
		repo: repo,
	}
}

func (s *Service) Publish(
	ctx context.Context,
	authorID uint64,
	title string,
	description string,
	videoURL string,
) (*Video, error) {
	// TODO 1：校验 authorID，不能为零
	if authorID == 0 {
		return nil, ErrInvalidAuthor
	}
	// TODO 2：用 strings.TrimSpace 处理
	// title、description、videoURL
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	videoURL = strings.TrimSpace(videoURL)
	// TODO 3：校验标题字符数
	// 不符合要求时返回 ErrInvalidTitle
	titleLength := utf8.RuneCountInString(title)
	if titleLength == 0 || titleLength > 100 {
		return nil, ErrInvalidTitle
	}
	// TODO 4：校验简介字符数
	// 不符合要求时返回 ErrInvalidDescription
	if utf8.RuneCountInString(description) > 1000 {
		return nil, ErrInvalidDescription
	}
	// TODO 5：校验视频地址
	// - 不能为空，最多 2048 字节
	// - url.Parse 不能报错
	if videoURL == "" || len(videoURL) > 2048 {
		return nil, ErrInvalidVideoURL
	}

	u, err := url.Parse(videoURL)
	if err != nil {
		return nil, ErrInvalidVideoURL
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return nil, ErrInvalidVideoURL
	}
	// - Scheme 必须为 http 或 https
	// - Hostname() 不能为空
	// 不符合要求时返回 ErrInvalidVideoURL
	// TODO 6：构造 Video，填入四个业务字段
	// ID 和 CreatedAt 由数据库/GORM 处理
	video := &Video{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		VideoURL:    videoURL,
	}
	// TODO 7：调用 s.repo.Create
	// 数据库失败时返回 nil 和原始错误
	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}
	// TODO 8：返回创建成功的 Video 和 nil
	return video, nil
}
