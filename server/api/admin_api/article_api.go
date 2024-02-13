package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 新建文章
func NewArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlCreate[model.Article](article)

	if err != nil {
		global.Logrus.Error(err)
		response.Fail("NewArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewArticle success", nil, ctx)
}

// 删除文章
func DeleteArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlDelete[model.Article, model.Article](article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteArticle success", nil, ctx)
}

// 更新文章
func UpdateArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = articleService.UpdateArticle(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateArticle error:"+err.Error(), nil, ctx)
		return
	}
	// 删除旧的缓存:默认的首页弹窗和自定义内容
	global.LocalCache.Delete(constant.CACHE_DEFAULT_ARTICLE)
	response.OK("UpdateArticle success", nil, ctx)

}

// 获取文章
func GetArticleList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	params.TableName = "article" //
	res, total, err := common_logic.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetArticle success", model.CommonDataResp{
		Total: total,
		Data:  res,
	}, ctx)
}
