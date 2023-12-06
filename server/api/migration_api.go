package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

func Migration(ctx *gin.Context) {
	var mig model.Migration
	err := ctx.ShouldBind(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	msg, err := service.Migration(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("Migration success:", msg, ctx)

}
