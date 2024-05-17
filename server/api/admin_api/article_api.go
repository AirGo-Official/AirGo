package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// NewArticle
// @Tags [admin api] article
// @Summary 新建文章
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Article true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/article/newArticle [post]
func NewArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.Article](article)

	if err != nil {
		global.Logrus.Error(err)
		response.Fail("NewArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewArticle success", nil, ctx)
}

// DeleteArticle
// @Tags [admin api] article
// @Summary 删除文章
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Article true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/article/deleteArticle [delete]
func DeleteArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Article, model.Article](article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteArticle success", nil, ctx)
}

// UpdateArticle
// @Tags [admin api] article
// @Summary 更新文章
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Article true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/article/updateArticle [post]
func UpdateArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminArticleSvc.UpdateArticle(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateArticle error:"+err.Error(), nil, ctx)
		return
	}
	// 删除旧的缓存:默认的首页弹窗和自定义内容
	global.LocalCache.Delete(constant.CACHE_DEFAULT_ARTICLE)
	response.OK("UpdateArticle success", nil, ctx)

}

// GetArticleList
// @Tags [admin api] article
// @Summary 获取文章列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/article/getArticleList [post]
func GetArticleList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	params.TableName = "article" //
	res, total, err := service.CommonSqlFindWithFieldParams(&params)
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
