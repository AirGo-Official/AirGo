package service

import (
	"AirGo/global"
	"AirGo/model"
)

// 根据uId查角色Ids
func FindRoleIdsByuId(uId int64) ([]int64, error) {
	var roles []model.UserAndRole
	err := global.DB.Model(&model.UserAndRole{}).Select("role_id").Where("user_id=?", uId).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	//角色id 数组
	var roleIds []int64
	for item := range roles {
		roleIds = append(roleIds, roles[item].RoleID)
	}
	return roleIds, nil
}

// 根据角色名的数组查询角色数组
func FindRoleIdsByRoleNameArr(nameArr []string) ([]model.Role, error) {
	var roles []model.Role
	//err := global.DB.Model(&model.Role{}).Select("id").Where("role_name in ?", nameArr).Find(&roles).Error
	err := global.DB.Model(&model.Role{}).Where("role_name in ?", nameArr).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil

}

// 分页获取角色列表
func GetRoleList(roleParams *model.PaginationParams) (*model.RolesWithTotal, error) {
	var roleArr model.RolesWithTotal
	var err error
	if roleParams.Search != "" {
		err = global.DB.Model(&model.Role{}).Count(&roleArr.Total).Where("role_name like ?", ("%" + roleParams.Search + "%")).Preload("Menus").Limit(int(roleParams.PageSize)).Offset((int(roleParams.PageNum) - 1) * int(roleParams.PageSize)).Find(&roleArr.RoleList).Error
	} else {
		err = global.DB.Model(&model.Role{}).Count(&roleArr.Total).Preload("Menus").Limit(int(roleParams.PageSize)).Offset((int(roleParams.PageNum) - 1) * int(roleParams.PageSize)).Find(&roleArr.RoleList).Error
	}
	return &roleArr, err

}

// 修改角色信息
func ModifyRoleInfo(roleInfo *model.Role) error {
	//先更新角色信息
	err := global.DB.Save(&roleInfo).Error
	if err != nil {
		return err
	}
	//再更新关联信息
	var routerSlice []model.DynamicRoute
	global.DB.Where("id in (?)", roleInfo.Nodes).Find(&routerSlice)
	err = global.DB.Model(&roleInfo).Association("Menus").Replace(&routerSlice)
	return err
}

// 新建角色
func AddRole(role *model.Role) error {
	//先查关联的动态路由
	//var routerSlice []model.SysDynamicRouter
	err := global.DB.Where("id in (?)", role.Nodes).Find(&role.Menus).Error
	if err != nil {
		return err
	}
	err = global.DB.Create(&role).Error
	return err
}

// 删除角色
func DelRole(id int64) error {
	var role model.Role
	//查询全部关联并删除
	err := global.DB.Where("id = ?", id).Preload("Menus").Find(&role).Error
	if err != nil {
		return err
	}
	global.DB.Model(&role).Association("Menus").Delete(role.Menus)
	//最后删除角色
	err = global.DB.Where("id = ?", id).Delete(&model.Role{}).Error
	return err
}

// 删除用户关联的角色组
func DeleteUserRoleGroup(user *model.User) error {
	return global.DB.Model(&model.User{ID: user.ID}).Association("RoleGroup").Replace(nil)
}
