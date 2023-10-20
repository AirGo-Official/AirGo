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
	//picUrl := ctx.Query("picUrl")
	//subject := ctx.Query("subject")

	var pic model.Gallery
	ctx.ShouldBind(&pic)

	if pic.PictureUrl == "" {
		response.Fail("上传图片链接参数错误", nil, ctx)
		return
	}
	if pic.Subject == "" {
		pic.Subject = time.Now().Format("2006-01-02 15:03:04")
	}
	//fmt.Printf("pic:%v", pic)
	uIDInt, _ := other_plugin.GetUserIDFromGinContext(ctx)
	pic.UserID = uIDInt
	err := service.CommonSqlCreate[model.Gallery](pic)
	//err := service.NewPictureUrl(uIDInt, picUrl, subject)
	if err != nil {
		response.Fail("上传图片链接错误", nil, ctx)
		return
	}
	response.OK("上传图片链接成功", nil, ctx)
}

// 获取图库列表
func GetPictureList(ctx *gin.Context) {
	var params model.PaginationParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error("获取图片列表错误：", err.Error())
		response.Fail("获取图片列表错误："+err.Error(), nil, ctx)
		return
	}
	var text string
	if params.Search != "" {
		text = "subject like" + " % " + params.Search + " % "
	}
	picList, _, err := service.CommonSqlFindWithPagination[model.Gallery, string, []model.Gallery](text, params)
	if err != nil {
		global.Logrus.Error("获取图片列表错误：", err.Error())
		response.Fail("获取图片列表错误："+err.Error(), nil, ctx)
		return
	}
	response.OK("获取图片列表成功", picList, ctx)
}
