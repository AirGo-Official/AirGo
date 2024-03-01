package user_logic

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
)

type Article struct{}

// 获取默认的首页弹窗和自定义内容
func (a *Article) GetDefaultArticle() (*model.CommonDataResp, error) {
	// 获取缓存
	cache, ok := global.LocalCache.Get(constant.CACHE_DEFAULT_ARTICLE)
	if ok {
		cache1 := cache.(model.CommonDataResp)
		return &cache1, nil
	}

	params := model.QueryParams{
		TableName: "article",
		FieldParamsList: []model.FieldParamsItem{
			{Operator: "",
				Field:          "id",
				Condition:      "=",
				ConditionValue: "1",
			},
			{
				Operator:       "OR",
				Field:          "id",
				Condition:      "=",
				ConditionValue: "2",
			},
			{
				Operator:       "AND",
				Field:          "status",
				Condition:      "=",
				ConditionValue: "1",
			},
		},
		Pagination: model.Pagination{
			PageNum:  1,
			PageSize: 2,
			OrderBy:  "id ASC",
		},
	}
	data, total, err := common_logic.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		return nil, err
	} else {
		//加入缓存
		articleList := data.([]model.Article)
		// 修改文章
		//strings.ReplaceAll(articleList[0].Content, "auto_replace_backend_url", global.Server.Website.BackendUrl) //替换backend url
		global.LocalCache.SetNoExpire(constant.CACHE_DEFAULT_ARTICLE, model.CommonDataResp{
			Total: total,
			Data:  articleList,
		})
		return &model.CommonDataResp{total, articleList}, nil
	}
}
