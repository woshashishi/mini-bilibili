package video

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// TODO：保存数据库连接，返回 Repository
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	video *Video,
) error {
	// TODO：
	// 1. 使用 WithContext(ctx) 传递请求上下文
	err := r.db.WithContext(ctx).Create(video).Error
	if err != nil {
		return err
	}
	// 2. 使用 Create 插入 video
	// 3. 返回数据库操作的 Error
	return nil
}
