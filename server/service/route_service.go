package service

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/format_plugin"
	"errors"
	"gorm.io/gorm"
)

// 角色Ids对应的route Ids
func GetRouteIdsByRoleIds(roleIds []int64) ([]int64, error) {
	var RoleAndMenuArr []model.RoleAndMenu
	if roleIds == nil {
		err := global.DB.Find(&RoleAndMenuArr).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := global.DB.Where("role_id in (?)", roleIds).Find(&RoleAndMenuArr).Error
		if err != nil {
			return nil, err
		}
	}

	var routeIds []int64
	for item := range RoleAndMenuArr {
		routeIds = append(routeIds, RoleAndMenuArr[item].DynamicRouteID)
	}
	//过滤重复
	routeIdsNew := format_plugin.ArrayDeduplication(routeIds)
	return routeIdsNew, nil

}

// 根据route Ids 查 route Slice
func GetRouteSliceByRouteIds(routeIds []int64) (*[]model.DynamicRoute, error) {
	var RouteArr []model.DynamicRoute

	if routeIds == nil {
		err := global.DB.Find(&RouteArr).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := global.DB.Where("id in (?)", routeIds).Find(&RouteArr).Error
		if err != nil {
			return nil, err
		}
	}

	return &RouteArr, nil
}

// 根据route Ids 查 route Node Slice
func GetRouteNodeByRouteIds(routeIds []int64) (*[]model.RouteNode, error) {
	var routeNodeSlice []model.RouteNode
	if routeIds == nil {
		err := global.DB.Model(model.DynamicRoute{}).Find(&routeNodeSlice).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := global.DB.Model(model.DynamicRoute{}).Where("id in (?)", routeIds).Find(&routeNodeSlice).Error
		if err != nil {
			return nil, err
		}
	}
	return &routeNodeSlice, nil
}

// 获取角色动态路由
func GetDynamicRoute(RouterSlice *[]model.DynamicRoute) *[]model.DynamicRoute {
	routeMap := make(map[int64][]model.DynamicRoute)
	for _, value := range *RouterSlice {
		routeMap[value.ParentID] = append(routeMap[value.ParentID], value)
	}
	routeSlice := routeMap[0]              //0为左侧顶级菜单,从左侧菜单开始，找到每一个菜单的子菜单
	for i := 0; i < len(routeSlice); i++ { //
		getChildrenRoute(&routeSlice[i], routeMap) //获取每一个一级菜单的子菜单
	}
	return &routeSlice
}

// 递归获取子路由
func getChildrenRoute(route *model.DynamicRoute, routeMap map[int64][]model.DynamicRoute) {
	route.Children = routeMap[route.ID]        //routeMap存的是所有:父子
	for i := 0; i < len(route.Children); i++ { //
		getChildrenRoute(&route.Children[i], routeMap) //
	}
}

// 角色的路由节点树
func GetRouteNodeTree(routeNodeSlice *[]model.RouteNode) *[]model.RouteNode {
	routeNodeMap := make(map[int64][]model.RouteNode)
	for _, value := range *routeNodeSlice {
		routeNodeMap[value.ParentID] = append(routeNodeMap[value.ParentID], value)
	}
	// 递归查询子菜单
	routeNode := routeNodeMap[0]
	for i := 0; i < len(routeNode); i++ {
		getChildrenRouteNodeTree(&routeNode[i], routeNodeMap)
	}
	return &routeNode
}

// 递归获取子路由节点树
func getChildrenRouteNodeTree(routeNode *model.RouteNode, routeNodeMap map[int64][]model.RouteNode) {
	routeNode.Children = routeNodeMap[routeNode.ID]
	for i := 0; i < len(routeNode.Children); i++ {
		getChildrenRouteNodeTree(&routeNode.Children[i], routeNodeMap)
	}
}

// 查询动态路由是否存在 by path
func NotExistDynamicRoute(route *model.DynamicRoute) bool {
	var dr model.DynamicRoute
	err := global.DB.Where(&model.DynamicRoute{Path: route.Path}).First(&dr).Error //注意Model Where 的区别
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// 查询动态路由 by meta.title
func FindDynamicRoute(route *model.DynamicRoute) ([]model.DynamicRoute, error) {
	var routeSlice []model.DynamicRoute
	err := global.DB.Where("title like ?", ("%" + route.Meta.Title + "%")).Find(&routeSlice).Error
	return routeSlice, err
}

// 新建动态路由
func NewDynamicRoute(route *model.DynamicRoute) error {
	err := global.DB.Create(&route).Error
	return err
}

// 删除动态路由
func DelDynamicRoute(route *model.DynamicRoute) error {
	//删除关联的路由
	err := global.DB.Where("dynamic_route_id = ?", route.ID).Delete(&model.RoleAndMenu{}).Error
	if err != nil {
		return err
	}
	//删除路由
	err = global.DB.Where(&route).Delete(&model.DynamicRoute{}).Error
	return err

}

// 修改动态路由
func UpdateDynamicRoute(route *model.DynamicRoute) error {
	//err := global.DB.Model(&route).Omit("router_id").Updates(&route).Error
	err := global.DB.Save(&route).Error
	return err
}
