package like

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// TODO：保存 db，返回 Repository
	return &Repository{
		db: db,
	}
}

func (r *Repository) Add(
	ctx context.Context,
	userID uint64,
	videoID uint64,
) error {
	// TODO 1：构造 Like
	like := &Like{
		UserID:  userID,
		VideoID: videoID,
	}
	// 只填 UserID、VideoID
	// 不要手动指定主键 ID

	// TODO 2：通过 WithContext(ctx).Create(...) 插入记录
	// 获取 Error
	err := r.db.WithContext(ctx).Create(like).Error
	// TODO 3：没有错误时，返回 nil
	if err == nil {
		return nil
	}
	// TODO 4：用 errors.As 提取 *mysql.MySQLError
	// 提取成功，且 Number == 1062：
	// 返回 nil，表示用户已经处于点赞状态
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		// mysqlErr.Number 是 MySQL 错误编号
		return nil
	}

	return err
}
func (r *Repository) Remove(
	ctx context.Context,
	userID uint64,
	videoID uint64,
) error {
	// TODO 1：使用 r.db.WithContext(ctx)

	// TODO 2：使用 Where 添加条件
	// "user_id = ? AND video_id = ?"
	// 参数依次为 userID、videoID
	err := r.db.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&Like{}).Error
	// TODO 3：调用 Delete(&Like{})
	// Like 模型用于确定操作 likes 表

	// TODO 4：返回操作的 Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) CountByVideoID(
	ctx context.Context,
	videoID uint64,
) (int64, error) {
	// TODO 1：声明 int64 类型的 count
	var count int64
	// TODO 2：构造查询：
	// WithContext(ctx)
	// Model(&Like{})
	// Where("video_id = ?", videoID)
	// Count(&count)
	// 获取 Error
	err := r.db.WithContext(ctx).Model(&Like{}).Where("video_id = ?", videoID).Count(&count).Error
	// TODO 3：失败时返回 0 和原始错误
	if err != nil {
		return 0, err
	}
	// TODO 4：成功时返回 count 和 nil
	return count, nil
}
func (r *Repository) Exists(
	ctx context.Context,
	userID uint64,
	videoID uint64,
) (bool, error) {
	// TODO 1：声明 int64 类型的 count
	var count int64
	// TODO 2：通过 WithContext(ctx) 构造查询
	// Model(&Like{})
	// Where 同时匹配 user_id 和 video_id
	// Count(&count)，获取 Error
	err := r.db.WithContext(ctx).Model(&Like{}).Where("user_id=?", userID).Where("video_id=?", videoID).Count(&count).Error
	// TODO 3：查询失败时，返回 false 和原始错误
	if err != nil {
		return false, err
	}
	// TODO 4：查询成功，返回 count > 0 和 nil
	return count > 0, nil
}
