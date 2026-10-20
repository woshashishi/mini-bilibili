package user

import "time"

type User struct {
	// TODO 1：
	// ID 使用 uint，并设置为主键。
	ID uint `gorm:"primaryKey"`
	// TODO 2：
	// Username：
	Username string `gorm:"uniqueIndex;type:varchar(32);not null"`
	// - 最大长度 32
	// - 不能为空
	// - 建立唯一索引

	// TODO 3：
	// PasswordHash：
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	// - 数据库列名使用 password_hash
	// - 最大长度 255
	// - 不能为空
	// - 返回 JSON 时不能暴露

	// TODO 4：
	// CreatedAt 使用 time.Time。
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// 让 GORM 在创建记录时自动填写。

}
