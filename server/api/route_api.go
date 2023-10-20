package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/other_plugin"
	"AirGo/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取角色动态路由
func GetRouteList(ctx *gin.Context) {
	uIdInt, _ := other_plugin.GetUserIDFromGinContext(ctx)
	//查询uId对应的角色
	roleIds, err := service.FindRoleIdsByuId(uIdInt)
	if err != nil {
		global.Logrus.Error("角色查询错误", err.Error())
		response.Fail("角色查询错误"+err.Error(), nil, ctx)
		return
	}
	// 角色Ids对应的route Ids
	routeIds, err := service.GetRouteIdsByRoleIds(roleIds)
	if err != nil {
		global.Logrus.Error("GetRouteIdsByRoleIds err", err)
		response.Fail("GetRouteIdsByRoleIds err"+err.Error(), nil, ctx)
		return
	}
	// 根据route Ids 查 route Slice
	routeSlice, err := service.GetRouteSliceByRouteIds(routeIds)
	if err != nil {
		global.Logrus.Error("GetRouteSliceByRouteIds err", err)
		response.Fail("GetRouteSliceByRouteIds err"+err.Error(), nil, ctx)
		return
	}
	// 获取角色动态路由
	route := service.GetDynamicRoute(routeSlice)
	response.OK("菜单获取成功", route, ctx)
}

// 获取全部角色动态路由
func GetAllRouteList(ctx *gin.Context) {
	// 根据route Ids 查 route Slice
	routeSlice, err := service.GetRouteSliceByRouteIds(nil)
	if err != nil {
		global.Logrus.Error("GetRouteSliceByRouteIds err", err)
		response.Fail("GetRouteSliceByRouteIds err"+err.Error(), nil, ctx)
		return
	}
	// 获取角色动态路由
	route := service.GetDynamicRoute(routeSlice)
	response.OK("全部角色动态路由获取成功", route, ctx)

}

// 前端编辑角色的时候显示全部菜单节点树
func GetAllRouteTree(ctx *gin.Context) {
	routeNodeSlice, err := service.GetRouteNodeByRouteIds(nil)
	if err != nil {
		global.Logrus.Error("GetRouteNodeByRouteIds err", err)
		response.Fail("GetRouteNodeByRouteIds err"+err.Error(), nil, ctx)
		return
	}
	routeNodeTree := service.GetRouteNodeTree(routeNodeSlice)

	response.OK("当前角色动态路由节点树获取成功", routeNodeTree, ctx)
}

// 前端编辑角色的时候显示当前角色的菜单tree
func GetRouteTree(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	// 角色Ids对应的route Ids
	var roleIds = []int64{roleId}
	routeIds, err := service.GetRouteIdsByRoleIds(roleIds) //空
	if err != nil {
		global.Logrus.Error("GetRouteIdsByRoleIds err", err)
		response.Fail("GetRouteIdsByRoleIds err"+err.Error(), nil, ctx)
		return
	}
	routeNodeSlice, err := service.GetRouteNodeByRouteIds(routeIds)
	if err != nil {
		global.Logrus.Error("GetRouteNodeByRouteIds err", err)
		response.Fail("GetRouteNodeByRouteIds err"+err.Error(), nil, ctx)
		return
	}
	routeNodeTree := service.GetRouteNodeTree(routeNodeSlice)
	response.OK("全部动态路由节点树获取成功", routeNodeTree, ctx)
}

// 新建动态路由
func NewDynamicRoute(ctx *gin.Context) {
	var route model.DynamicRoute
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error("新建动态路由参数错误", err.Error())
		response.Fail("新建动态路由参数错误"+err.Error(), nil, ctx)
		return
	}
	route.ID = 0
	// 查询动态路由是否存在
	notExist := service.NotExistDynamicRoute(&route)
	if !notExist {
		response.Fail("动态路由已存在", nil, ctx)
		return
	}
	err = service.NewDynamicRoute(&route)
	if err != nil {
		global.Logrus.Error("新建动态路由参数错误", err.Error())
		response.Fail("新建动态路由错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("新建动态路由成功", nil, ctx)

}

// 删除动态路由
func DelDynamicRoute(ctx *gin.Context) {
	var route model.DynamicRoute
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error("删除动态路由参数错误", err.Error())
		response.Fail("动态路由参数错误"+err.Error(), nil, ctx)
		return
	}
	// 查询动态路由是否存在
	notExist := service.NotExistDynamicRoute(&route)
	if notExist {
		response.Fail("动态路由不存在", nil, ctx)
		return
	}
	err = service.DelDynamicRoute(&route)
	if err != nil {
		global.Logrus.Error("删除动态路由参数错误", err.Error())
		response.Fail("动态路由错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("删除动态路由成功", nil, ctx)

}

// 修改动态路由
func UpdateDynamicRoute(ctx *gin.Context) {
	var route model.DynamicRoute
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error("修改动态路由参数错误", err.Error())
		response.Fail("动态路由参数错误"+err.Error(), nil, ctx)
		return
	}
	// 查询动态路由是否存在
	notExist := service.NotExistDynamicRoute(&route)
	if notExist {
		response.Fail("动态路由不存在", nil, ctx)
		return
	}

	err = service.UpdateDynamicRoute(&route)
	if err != nil {
		global.Logrus.Error("修改动态路由参数错误", err.Error())
		response.Fail("动态路由错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改动态路由成功", nil, ctx)

}

// 查询单条动态路由 by meta.title
func FindDynamicRoute(ctx *gin.Context) {
	var route model.DynamicRoute
	err := ctx.ShouldBind(&route)
	if err != nil {
		global.Logrus.Error("查询动态路由参数错误", err.Error())
		response.Fail("单条动态路由参数错误"+err.Error(), nil, ctx)
		return
	}
	res, err := service.FindDynamicRoute(&route)
	if err != nil {
		global.Logrus.Error("查询动态路由参数错误", err.Error())
		response.Fail("单条动态路由错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("查询单条动态路由成功", res, ctx)

}
