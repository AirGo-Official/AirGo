package user_api

import (
	"github.com/AirGo-Official/AirGo/constant"
	"github.com/AirGo-Official/AirGo/global"
	"github.com/AirGo-Official/AirGo/model"
	service "github.com/AirGo-Official/AirGo/service"
	"github.com/AirGo-Official/AirGo/utils/response"
	"github.com/gin-gonic/gin"
)

// GetArticleList
// @Tags [customer api] article
// @Summary 获取文章列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/article/getArticleList [post]
func GetArticleList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	params.TableName = "article" //查询article表
	res, total, err := service.CommonSqlFindWithFieldParams(&params)
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
