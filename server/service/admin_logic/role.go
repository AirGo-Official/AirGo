package admin_logic

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Role struct{}

var roleService *Role

// 根据uId查角色Ids
func (r *Role) FindRoleIdsByuId(uId int64) ([]int64, error) {
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
func (r *Role) FindRoleIdsByRoleNameArr(nameArr []string) ([]model.Role, error) {
	var roles []model.Role
	err := global.DB.Model(&model.Role{}).Where("role_name in ?", nameArr).Find(&roles).Error
	return roles, err

}

// 获取角色列表
func (r *Role) GetRoleList() (*model.CommonDataResp, error) {
	var roleArr []model.Role
	var total int64
	var err error
	err = global.DB.
		Model(&model.Role{}).
		Count(&total).
		Preload("Menus").
		Find(&roleArr).Error
	return &model.CommonDataResp{total, roleArr}, err
}

// 修改角色信息
func (r *Role) UpdateRole(roleParams *model.Role) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//更新关联
		err := tx.Model(&roleParams).Association("Menus").Replace(&roleParams.Menus)
		if err != nil {
			return err
		}
		//更新角色信息
		return tx.Save(&roleParams).Error
	})
}

// 新建角色
func (r *Role) NewRole(roleParams *model.Role) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&roleParams).Error
	})
}

// 删除角色
func (r *Role) DelRole(roleParams *model.Role) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除多对多关联的菜单
		err := tx.Model(&roleParams).Association("Menus").Replace(nil)
		if err != nil {
			return err
		}
		// 删除多对多关联的用户
		err = tx.Model(&roleParams).Association("UserGroup").Replace(nil)
		if err != nil {
			return err
		}
		//最后删除角色
		return tx.Where(&model.Role{ID: roleParams.ID}).Delete(&model.Role{}).Error
	})
}

// 删除用户关联的角色组
func (r *Role) DeleteUserRoleGroup(userParams *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.User{ID: userParams.ID}).Association("RoleGroup").Replace(nil)
	})
}
