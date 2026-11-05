package video

import "time"

type Video struct {
	// TODO：主键、自增，JSON 字段名为 id
	ID uint64 `gorm:"primaryKey" json:"id"`

	// TODO：非空、普通索引，JSON 字段名为 author_id
	AuthorID uint64 `gorm:"not null;index" json:"author_id"`

	// TODO：varchar(100)、非空，JSON 字段名为 title
	Title string `gorm:"not null;type:varchar(100)" json:"title"`

	// TODO：varchar(1000)，JSON 字段名为 description
	Description string `gorm:"type:varchar(1000)" json:"description"`

	// TODO：varchar(2048)、非空，JSON 字段名为 video_url
	VideoURL string `gorm:"type:varchar(2048);not null" json:"video_url"`

	// TODO：自动记录创建时间、非空，JSON 字段名为 created_at
	CreatedAt time.Time `gorm:"autoCreateTime;not null;" json:"created_at"`
}
