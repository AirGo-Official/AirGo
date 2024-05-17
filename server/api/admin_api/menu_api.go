package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetAllMenuList
// @Tags [admin api] menu
// @Summary 获取全部菜单列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/menu/getAllMenuList [get]
func GetAllMenuList(ctx *gin.Context) {
	routeList, err := service.AdminMenuSvc.GetMenuList()
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenusByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	route := service.AdminMenuSvc.GetMenus(routeList)
	response.OK("GetAllMenuList success", route, ctx)
}

// NewMenu
// @Tags [admin api] menu
// @Summary 新建菜单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Menu true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/menu/newMenu [post]
func NewMenu(ctx *gin.Context) {
	var menu model.Menu
	err := ctx.ShouldBind(&menu)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	menu.ID = 0
	err = service.AdminMenuSvc.NewMenu(&menu)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewMenu success", nil, ctx)
}

// DelMenu
// @Tags [admin api] menu
// @Summary 删除菜单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Menu true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/menu/delMenu [delete]
func DelMenu(ctx *gin.Context) {
	var route model.Menu
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminMenuSvc.DelMenu(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DelMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DelMenu success", nil, ctx)

}

// UpdateMenu
// @Tags [admin api] menu
// @Summary 修改菜单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Menu true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/menu/updateMenu [post]
func UpdateMenu(ctx *gin.Context) {
	var route model.Menu
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminMenuSvc.UpdateMenu(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateMenu success", nil, ctx)
}
