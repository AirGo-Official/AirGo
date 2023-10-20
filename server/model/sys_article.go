package model

import "time"

type Article struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"-" gorm:"index"`
	ID int64 `json:"id"      gorm:"primary_key"`

	Type         string `json:"type"         gorm:"comment:文章类型，home=首页自定义显示内容，dialog=首页弹窗，notice=公告，knowledge=知识库"`
	Status       bool   `json:"status"       gorm:"comment:是否启用"`
	Title        string `json:"title"        gorm:"comment:文章标题"`
	Introduction string `json:"introduction" gorm:"comment:文章简介"`
	Content      string `json:"content"      gorm:"comment:文章内容;size:30000"`
}

type ArticleWithTotal struct {
	Total       int64     `json:"total"`
	ArticleList []Article `json:"article_list"`
}
