package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

var ms admin_logic.Migration

func Migration(ctx *gin.Context) {
	var mig model.Migration
	err := ctx.ShouldBind(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	msg, err := ms.Migration(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("Migration success:", msg, ctx)

}
