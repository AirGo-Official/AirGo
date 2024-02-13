package admin_logic

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Article struct{}

// 更新文章
func (a *Article) UpdateArticle(article *model.Article) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&article).Error
	})
}
