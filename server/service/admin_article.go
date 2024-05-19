package service

import (
	"github.com/AirGo-Official/AirGo/constant"
	"github.com/AirGo-Official/AirGo/global"
	"github.com/AirGo-Official/AirGo/model"
	"gorm.io/gorm"
)

type AdminArticle struct{}

var AdminArticleSvc *AdminArticle

// 更新文章
func (a *AdminArticle) UpdateArticle(article *model.Article) error {
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
