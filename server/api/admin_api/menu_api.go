package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取全部动态路由
func GetAllMenuList(ctx *gin.Context) {
	routeList, err := menuService.GetMenuList()
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenusByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	route := menuService.GetMenus(routeList)
	response.OK("GetAllMenuList success", route, ctx)
}

// 新建动态路由
func NewMenu(ctx *gin.Context) {
	var menu model.Menu
	err := ctx.ShouldBind(&menu)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	menu.ID = 0
	err = menuService.NewMenu(&menu)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewMenu success", nil, ctx)
}

// 删除动态路由
func DelMenu(ctx *gin.Context) {
	var route model.Menu
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = menuService.DelMenu(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DelMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DelMenu success", nil, ctx)

}

// 修改动态路由
func UpdateMenu(ctx *gin.Context) {
	var route model.Menu
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = menuService.UpdateMenu(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateMenu error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateMenu success", nil, ctx)
}
