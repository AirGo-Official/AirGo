package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetArticleList 获取文章
func GetArticleList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	params.TableName = "article" //查询article表
	res, total, err := common_logic.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetArticleList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetArticleList success", model.CommonDataResp{
		Total: total,
		Data:  res,
	}, ctx)
}

// 获取默认的首页弹窗和自定义内容
func GetDefaultArticleList(ctx *gin.Context) {
	data, err := articleService.GetDefaultArticle()
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetDefaultArticleList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetDefaultArticleList success", data, ctx)
}
