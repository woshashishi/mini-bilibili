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
func (r *Repository) FindByID(
	ctx context.Context,
	id uint64,
) (*Video, error) {
	// TODO 1：声明一个 Video 变量
	video := &Video{
		ID: id,
	}
	// TODO 2：使用 r.db.WithContext(ctx)
	// 按 id 查询第一条记录
	// 使用参数占位符，条件为 "id = ?"
	err := r.db.WithContext(ctx).Where("id=?", id).First(video).Error
	// TODO 3：查询失败时，返回 nil 和数据库错误
	// 包括 gorm.ErrRecordNotFound，暂时不要转换它
	if err != nil {
		return nil, err
	}
	// TODO 4：查询成功时，返回 Video 指针和 nil
	return video, nil
}
func (r *Repository) ListBeforeID(
	ctx context.Context,
	cursor uint64,
	limit int,
) ([]Video, error) {
	// TODO 1：创建空的 Video 切片
	// 建议使用 make([]Video, 0)
	var video = make([]Video, 0)
	// TODO 2：从 r.db.WithContext(ctx) 创建查询对象
	query := r.db.WithContext(ctx)
	// TODO 3：cursor > 0 时添加 "id < ?" 条件
	// 注意保存 Where 返回的查询对象
	if cursor > 0 {
		query = query.Where("id <?", cursor)
	}
	// TODO 4：
	// 按 "id DESC" 排序
	// 使用 Limit(limit) 限制数量
	// 使用 Find 将结果填入切片
	err := query.Order("id DESC").Limit(limit).Find(&video).Error
	// TODO 5：查询失败时返回 nil 和原始错误
	if err != nil {
		return nil, err
	}
	// TODO 6：查询成功时返回切片和 nil
	return video, nil
}
