package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// TODO：返回 Repository，并保存 db

	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	// TODO：
	// 1. 声明 User 变量
	var user User
	// 2. 使用 r.db.WithContext(ctx)

	// 3. 按 username 查询
	// 4. 查询成功返回 &user
	err := r.db.WithContext(ctx).Where("username=?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
	// 5. 查询失败返回 nil 和 error

}

func (r *Repository) Create(
	ctx context.Context,
	user *User,
) error {
	// TODO：
	// 使用 WithContext(ctx).Create(...)
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
