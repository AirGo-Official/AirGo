package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type AdminMenu struct{}

var AdminMenuSvc *AdminMenu

// 新建菜单
func (m *AdminMenu) NewMenu(menu *model.Menu) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&menu).Error
	})
}

// 删除菜单
func (m *AdminMenu) DelMenu(menu *model.Menu) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//删除与角色的关联
		err := tx.Model(&menu).Association("Roles").Replace(nil)
		if err != nil {
			return err
		}
		//删除路由
		return tx.Delete(&model.Menu{}, menu.ID).Error
	})
}

// 修改菜单
func (m *AdminMenu) UpdateMenu(menu *model.Menu) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&menu).Error
	})
}

// 查 menu list
func (m *AdminMenu) GetMenuList() (*[]model.Menu, error) {
	var menuArr []model.Menu
	err := global.DB.Find(&menuArr).Error
	if err != nil {
		return nil, err
	}
	return &menuArr, nil
}

// 处理menu list，获取角色菜单
func (m *AdminMenu) GetMenus(menuList *[]model.Menu) *[]model.Menu {
	menuMap := make(map[int64][]model.Menu)
	for _, value := range *menuList {
		menuMap[value.ParentID] = append(menuMap[value.ParentID], value)
	}
	newMenuList := menuMap[0]               //0为左侧顶级菜单,从左侧菜单开始，找到每一个菜单的子菜单
	for i := 0; i < len(newMenuList); i++ { //
		m.getChildrenMenus(&newMenuList[i], menuMap) //获取每一个一级菜单的子菜单
	}
	return &newMenuList
}

// 递归获取子路由
func (m *AdminMenu) getChildrenMenus(menu *model.Menu, menuMap map[int64][]model.Menu) {
	menu.Children = menuMap[menu.ID]          //menuMap存的是所有:父子
	for i := 0; i < len(menu.Children); i++ { //
		m.getChildrenMenus(&menu.Children[i], menuMap) //
	}
}
