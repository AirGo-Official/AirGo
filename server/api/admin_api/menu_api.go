package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
	"strconv"
)

// 获取全部动态路由
func GetAllMenuList(ctx *gin.Context) {
	// 根据route Ids 查 route Slice
	routeSlice, err := menuService.GetMenusByMenuIds(nil)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenusByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	route := menuService.GetMenus(routeSlice)
	response.OK("GetAllMenuList success", route, ctx)

}

// 前端编辑角色的时候显示全部菜单节点树
func GetAllMenuTree(ctx *gin.Context) {
	routeNodeSlice, err := menuService.GetMenuNodeByMenuIds(nil)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenuNodeByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	routeNodeTree := menuService.GetMenusNodeTree(routeNodeSlice)
	response.OK("GetAllMenuTree success", routeNodeTree, ctx)
}

// 前端编辑角色的时候显示当前角色的菜单tree
func GetMenuTree(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	// 角色Ids对应的route Ids
	var roleIds = []int64{roleId}
	routeIds, err := menuService.GetMenuIdsByRoleIds(roleIds)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenuIdsByRoleIds error:"+err.Error(), nil, ctx)
		return
	}
	routeNodeSlice, err := menuService.GetMenuNodeByMenuIds(routeIds)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenuNodeByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	routeNodeTree := menuService.GetMenusNodeTree(routeNodeSlice)
	response.OK("GetMenuTree success", routeNodeTree, ctx)
}

// 新建动态路由
func NewMenu(ctx *gin.Context) {
	var route model.Menu
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	route.ID = 0
	err = menuService.NewMenu(&route)
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
