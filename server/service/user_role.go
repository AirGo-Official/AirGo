package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
)

type Role struct {
}

var RoleSvc *Role

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
