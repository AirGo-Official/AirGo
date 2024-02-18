package admin_logic

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Article struct{}

var articleService *Article

// 更新文章
func (a *Article) UpdateArticle(article *model.Article) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&article).Error
	})
	if err != nil {
		return err
	}
	//删除缓存
	global.LocalCache.Delete(constant.CACHE_DEFAULT_ARTICLE)
	return nil
}
