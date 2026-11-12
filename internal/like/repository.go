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
