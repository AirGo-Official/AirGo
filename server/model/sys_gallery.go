package model

import "time"

type Gallery struct {
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" gorm:"index"`
	ID         int64      `json:"id" gorm:"primary_key"`
	Subject    string     `json:"subject"       gorm:"comment:图片标题"`
	PictureUrl string     `json:"picture_url"   gorm:"comment:图片链接"`
	UserID     int64      `json:"user_id"       gorm:"comment:用户ID"`
}
