package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/other_plugin"
)

type Menu struct{}

var MenuSvc *Menu

// 角色Ids对应的menu Ids
func (m *Menu) GetMenuIdsByRoleIds(roleIds []int64) ([]int64, error) {
	var RoleAndMenuArr []model.RoleAndMenu
	err := global.DB.Where("role_id in (?)", roleIds).Find(&RoleAndMenuArr).Error
	if err != nil {
		return nil, err
	}

	var menuIds []int64
	for item := range RoleAndMenuArr {
		menuIds = append(menuIds, RoleAndMenuArr[item].MenuID)
	}
	//过滤重复
	return other_plugin.ArrayDeduplication(menuIds), nil
}

// 根据menu Ids 查 menu Slice
func (m *Menu) GetMenusByMenuIds(menuIds []int64) (*[]model.Menu, error) {
	var menuArr []model.Menu
	err := global.DB.Where("id in (?)", menuIds).Find(&menuArr).Error
	if err != nil {
		return nil, err
	}
	return &menuArr, nil
}

// 获取角色菜单
func (m *Menu) GetMenus(menuSlice *[]model.Menu) *[]model.Menu {
	menuMap := make(map[int64][]model.Menu)
	for _, value := range *menuSlice {
		menuMap[value.ParentID] = append(menuMap[value.ParentID], value)
	}
	newMenuSlice := menuMap[0]               //0为左侧顶级菜单,从左侧菜单开始，找到每一个菜单的子菜单
	for i := 0; i < len(newMenuSlice); i++ { //
		m.getChildrenMenus(&newMenuSlice[i], menuMap) //获取每一个一级菜单的子菜单
	}
	return &newMenuSlice
}

// 递归获取子菜单
func (m *Menu) getChildrenMenus(menu *model.Menu, menuMap map[int64][]model.Menu) {
	menu.Children = menuMap[menu.ID]          //menuMap存的是所有:父子
	for i := 0; i < len(menu.Children); i++ { //
		m.getChildrenMenus(&menu.Children[i], menuMap) //
	}
}
