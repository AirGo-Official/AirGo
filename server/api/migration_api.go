package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
)

func Migration(ctx *gin.Context) {
	var mig model.Migration
	err := ctx.ShouldBind(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	service.Show(mig)
	msg, err := service.Migration(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("Migration success:", msg, ctx)

}
