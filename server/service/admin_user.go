package service

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
	"time"
)

type AdminUser struct{}

var AdminUserSvc *AdminUser

// 新建用户
func (u *AdminUser) NewUser(userParams model.User) error {
	//判断是否存在
	var userQuery model.User
	err := global.DB.Where(&model.User{UserName: userParams.UserName}).First(&userQuery).Error
	if err == nil {
		return errors.New("User already exists")
	} else {
		if userParams.Avatar == "" {
			point := strings.Index(userParams.UserName, "@qq")
			if point != -1 {
				userParams.Avatar = fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%s&s=100", userParams.UserName[0:point])
			} else {
				userParams.Avatar = fmt.Sprintf("https://api.multiavatar.com/%s.svg", userParams.UserName)
			}
		}
		userParams.Password = encrypt_plugin.BcryptEncode(userParams.Password)
		userParams.InvitationCode = encrypt_plugin.RandomString(8)
		return global.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&userParams).Error
		})
	}
}

// 查用户
func (u *AdminUser) FirstUser(userParams *model.User) (*model.User, error) {
	var userQuery model.User
	err := global.DB.Where(&userParams).First(&userQuery).Error
	return &userQuery, err
}

// 保存用户
func (u *AdminUser) SaveUser(userParams *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&userParams).Error
	})
}
func (u *AdminUser) UpdateUser(userParams *model.User, values map[string]any) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.User{}).Where(&userParams).Updates(values).Error
	})
}

// 删除用户
func (u *AdminUser) DeleteUser(userParams *model.User) error {
	//删除user，删除全部has many，has one，many to many，不删除belongs to
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Select(clause.Associations).Delete(&userParams).Error
		if err != nil {
			return err
		}
		return tx.Where("user_id = ?", userParams.ID).Delete(&model.CustomerService{}).Error
	})
}

func (u *AdminUser) DeleteUserCacheTokenByID(userParams *model.User) {
	global.LocalCache.Delete(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, userParams.ID))
}

// 获取用户列表
func (u *AdminUser) GetUserlist(params *model.QueryParams) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var userList []model.User
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	err := global.DB.Model(&model.User{}).Count(&data.Total).Where(dataSql).Preload("RoleGroup").Find(&userList).Error
	if err != nil {
		return nil, err
	}
	data.Data = userList
	return &data, nil
}

// 重置管理员密码
func (u *AdminUser) ResetAdminPassword() {
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

func (u *AdminUser) UserSummary(params *model.QueryParams) (*[]model.UserSummary, error) {
	//处理查询时间
	var startTime, endTime time.Time
	startTime, err := time.Parse("2006-01-02 15:04:05", params.FieldParamsList[0].ConditionValue)
	if err != nil {
		return nil, err
	}
	endTime, _ = time.Parse("2006-01-02 15:04:05", params.FieldParamsList[1].ConditionValue)
	if err != nil {
		return nil, err
	}
	const (
		sql1 = `SELECT
DATE(created_at) as date,
COUNT(id) AS register_total
	`
		sql2 = "FROM `user`"
	)
	sql3 := fmt.Sprintf(" WHERE created_at > '%s' AND created_at < '%s'  GROUP BY date", startTime, endTime)
	var userSummary []model.UserSummary
	err = global.DB.
		Raw(sql1 + sql2 + sql3).
		Scan(&userSummary).
		Error
	//fmt.Println("result:", userSummary, "err:", err)
	return &userSummary, err
}
