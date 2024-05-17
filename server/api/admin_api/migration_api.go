package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"

	"github.com/ppoonk/AirGo/utils/response"
)

// Migration
// @Tags [admin api] migration
// @Summary 数据迁移
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Migration true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/migration/migrationData [post]
func Migration(ctx *gin.Context) {
	var mig model.Migration
	err := ctx.ShouldBind(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	msg, err := service.AdminMigrationSvc.Migration(&mig)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("Migration error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("Migration success:", msg, ctx)

}
