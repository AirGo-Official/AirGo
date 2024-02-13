package user_logic

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
)

type Article struct{}

func (a *Article) GetArticle(params *model.QueryParams) (any, int64, error) {
	return common_logic.CommonSqlFindWithFieldParams(params)
}

// 获取默认的首页弹窗和自定义内容
func (a *Article) GetDefaultArticle() (any, int64, error) {
	// 获取缓存
	cache, ok := global.LocalCache.Get(constant.CACHE_DEFAULT_ARTICLE)
	if ok {
		cache1 := cache.(model.CommonDataResp)
		return cache1.Data, cache1.Total, nil
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
		},
		Pagination: model.Pagination{
			PageNum:  1,
			PageSize: 2,
			OrderBy:  "id ASC",
		},
	}
	data, total, err := a.GetArticle(&params)
	if err != nil {
		return nil, 0, err
	} else {
		//加入缓存
		// TODO 修改文章
		global.LocalCache.SetNoExpire(constant.CACHE_DEFAULT_ARTICLE, model.CommonDataResp{
			Total: total,
			Data:  data,
		})
		return data, total, nil
	}

}
