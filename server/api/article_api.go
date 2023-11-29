package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取文章
func GetArticle(ctx *gin.Context) {
	var params model.PaginationParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetArticle error:"+err.Error(), nil, ctx)
		return
	}
	res, total, err := service.CommonSqlFindWithPagination[model.Article, string, []model.Article](params.Search, params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetArticle error:"+err.Error(), nil, ctx)
		return
	}
	var list = model.ArticleWithTotal{
		Total:       total,
		ArticleList: res,
	}
	response.OK("GetArticle success", list, ctx)
}

// 新建文章
func NewArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("NewArticle error:"+err.Error(), nil, ctx)
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

// 删除文章
func DeleteArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteArticle error:"+err.Error(), nil, ctx)
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

// 更新文章
func UpdateArticle(ctx *gin.Context) {
	var article model.Article
	err := ctx.ShouldBind(&article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateArticle error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Article](article)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateArticle error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateArticle success", nil, ctx)

}
