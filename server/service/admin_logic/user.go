package admin_logic

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

type User struct{}

var userService *User

// 新建用户
func (u *User) NewUser(user model.User) error {
	//判断是否存在
	var userQuery model.User
	err := global.DB.Where(&model.User{UserName: user.UserName}).First(&userQuery).Error
	if err == nil {
		return errors.New("User already exists")
	} else {
		//处理角色, 前端传的是角色名字，
		// todo 改为id
		var roleArr []string
		for _, v := range user.RoleGroup {
			roleArr = append(roleArr, v.RoleName)
		}
		roles, err := roleService.FindRoleIdsByRoleNameArr(roleArr)
		if err != nil {
			return err
		}
		user.RoleGroup = roles
		user.Password = encrypt_plugin.BcryptEncode(user.Password)
		user.InvitationCode = encrypt_plugin.RandomString(8)
		return global.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&u).Error
		})
	}
}

// 查用户
func (u *User) FirstUser(user *model.User) (*model.User, error) {
	var userQuery model.User
	err := global.DB.Where(&user).First(&userQuery).Error
	return &userQuery, err
}

// 保存用户
func (u *User) SaveUser(user *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&user).Error
	})
}

// 删除用户
func (u *User) DeleteUserOld(user *model.User) error {
	//删除关联的角色组
	err := roleService.DeleteUserRoleGroup(user)
	if err != nil {
		return err
	}
	//删除关联的订单
	err = orderService.DeleteUserAllOrder(user)
	if err != nil {
		return err
	}
	//删除用户
	return global.DB.Delete(&u).Error
}

// 删除用户
func (u *User) DeleteUser(userParams *model.User) error {
	return global.DB.Select(clause.Associations).Delete(&userParams).Error
}

func (u *User) DeleteUserCacheTokenByID(user *model.User) {
	global.LocalCache.Delete(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, user.ID))
}

// 获取用户列表
func (u *User) GetUserlist(params *model.QueryParams) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var userList []model.User
	_, dataSql := common_logic.CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	err := global.DB.Model(&model.User{}).Count(&data.Total).Where(dataSql).Preload("RoleGroup").Find(&userList).Error
	if err != nil {
		return nil, err
	}
	data.Data = userList
	return &data, nil
}

// 重置管理员密码
func (u *User) ResetAdminPassword() {
	//判断是否连接数据库
	if global.DB == nil {
		fmt.Println("未连接数据库，退出...")
		return
	}
	var name, psw, msg string
	fmt.Println("【重置管理员账户密码】\n请输入管理员邮箱(长度4～40)：")
	if _, err := fmt.Scanln(&name); err != nil {
		fmt.Println("输入出错了，已退出", err.Error())
		return
	}
	fmt.Println("请输入管理员密码(长度4～20)：")
	if _, err := fmt.Scanln(&psw); err != nil {
		fmt.Println("输入出错了，已退出", err.Error())
		return
	}
	fmt.Printf("账户：%s，密码：%s，是否重置？\n输入y重置，输入n退出\n", name, psw)
	if _, err := fmt.Scanln(&msg); err != nil {
		fmt.Println("输入出错了，已退出", err.Error())
		return
	}
	switch msg {
	case "y":
		fmt.Println("正在重置管理员账户密码...")
		var user model.User
		global.DB.First(&user)
		user.UserName = name
		user.Password = encrypt_plugin.BcryptEncode(psw)
		u.SaveUser(&user)
		fmt.Println("完成...")

	case "n":
		os.Exit(0)
	}

}
