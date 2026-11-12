package like

import "time"

type Like struct {
	// TODO：主键、自增
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// TODO：非空
	// 与 VideoID 组成同一个联合唯一索引
	UserID uint64 `gorm:"not null;uniqueIndex:idx_user_video"`

	// TODO：非空
	// 与 UserID 组成同一个联合唯一索引
	// 另外增加一个普通索引，便于按视频统计点赞数
	VideoID uint64 `gorm:"not null;uniqueIndex:idx_user_video;index"`

	// TODO：自动创建时间、非空
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
