package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/other_plugin"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
	"time"
)

// 上传图片链接
func NewPictureUrl(ctx *gin.Context) {
	var pic model.Gallery
	ctx.ShouldBind(&pic)

	if pic.PictureUrl == "" {
		response.Fail("NewPictureUrl error", nil, ctx)
		return
	}
	if pic.Subject == "" {
		pic.Subject = time.Now().Format("2006-01-02 15:03:04")
	}
	uIDInt, _ := other_plugin.GetUserIDFromGinContext(ctx)
	pic.UserID = uIDInt
	err := service.CommonSqlCreate[model.Gallery](pic)
	if err != nil {
		response.Fail("NewPictureUrl error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewPictureUrl success", nil, ctx)
}

// 获取图库列表
func GetPictureList(ctx *gin.Context) {
	var params model.PaginationParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetPictureList error:"+err.Error(), nil, ctx)
		return
	}
	var text string
	if params.Search != "" {
		text = "subject like" + " % " + params.Search + " % "
	}
	picList, _, err := service.CommonSqlFindWithPagination[model.Gallery, string, []model.Gallery](text, params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetPictureList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetPictureList success", picList, ctx)
}
