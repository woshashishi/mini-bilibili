package video

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidAuthor      = errors.New("invalid author")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidVideoURL    = errors.New("invalid video url")
	ErrInvalidVideoID     = errors.New("invalid video id")
	ErrVideoNotFound      = errors.New("video not found")
	ErrInvalidLimit       = errors.New("invalid limit")
)

type Service struct {
	repo *Repository
}

type FeedResult struct {
	Items      []Video
	HasMore    bool
	NextCursor uint64
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
func (s *Service) GetByID(
	ctx context.Context,
	id uint64,
) (*Video, error) {
	// TODO 1：ID 为零时
	// 返回 nil, ErrInvalidVideoID
	if id == 0 {
		return nil, ErrInvalidVideoID
	}
	// TODO 2：调用 s.repo.FindByID
	video, err := s.repo.FindByID(ctx, id)
	// TODO 3：使用 errors.Is 判断 gorm.ErrRecordNotFound
	// 转换为 ErrVideoNotFound

	// TODO 4：其他数据库错误原样返回
	if err != nil {
		ok := errors.Is(err, gorm.ErrRecordNotFound)
		if ok {
			return nil, ErrVideoNotFound
		}
		return nil, err
	}
	// TODO 5：成功返回视频和 nil
	return video, nil
}

func (s *Service) Feed(
	ctx context.Context,
	cursor uint64,
	limit int,
) (*FeedResult, error) {
	feedResult := &FeedResult{}
	// TODO 1：limit 为零时，设置默认值 10
	if limit == 0 {
		limit = 10
	}
	// TODO 2：limit 小于 1 或大于 50 时
	// 返回 nil, ErrInvalidLimit
	if limit < 1 || limit > 50 {
		return nil, ErrInvalidLimit
	}
	// TODO 3：调用 ListBeforeID
	// 注意：查询数量是 limit + 1
	// 数据库失败时返回 nil 和原始错误
	video, err := s.repo.ListBeforeID(ctx, cursor, limit+1)
	if err != nil {
		return nil, err
	}
	// TODO 4：判断查询结果是否超过 limit
	// 超过时：
	// - HasMore 为 true
	// - 将列表截取为前 limit 条
	//
	// 不超过时：
	// - HasMore 为 false
	// - 保留全部结果
	if len(video) > limit {
		feedResult.HasMore = true
		video = video[:limit]
	}
	// TODO 5：计算下一游标
	// 默认是 0
	// 只有 HasMore 为 true 时，
	// 才取截取后列表最后一条视频的 ID
	feedResult.NextCursor = 0
	if feedResult.HasMore == true {
		feedResult.NextCursor = video[len(video)-1].ID
	}
	// TODO 6：构造并返回 FeedResult
	feedResult.Items = video
	return feedResult, nil
}
