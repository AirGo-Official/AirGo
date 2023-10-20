package service

import (
	"AirGo/global"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"strings"
	"time"

	"AirGo/model"
	encrypt_plugin "AirGo/utils/encrypt_plugin"
	"errors"
	uuid "github.com/satori/go.uuid"
)

// 注册,角色对应如下
// {ID: 1, RoleName: "admin", Description: "超级管理员"},
// {ID: 2, RoleName: "普通用户", Description: "普通用户"},
func Register(u *model.User) error {
	//判断是否存在
	var user model.User
	err := global.DB.Where(&model.User{UserName: u.UserName}).First(&user).Error
	if err == nil {
		return errors.New("用户已存在")
	} else if err == gorm.ErrRecordNotFound {
		var newUser = model.User{
			UUID:           uuid.NewV4(),
			UserName:       u.UserName,
			NickName:       u.UserName,
			Password:       encrypt_plugin.BcryptEncode(u.Password),
			RoleGroup:      []model.Role{{ID: 2}}, //默认角色
			InvitationCode: encrypt_plugin.RandomString(8),
			ReferrerCode:   u.ReferrerCode,
		}
		return CreateUser(NewUserSubscribe(&newUser))
	} else {
		return err
	}
}

// 新建用户
func NewUser(u model.User) error {
	//判断是否存在
	var user model.User
	err := global.DB.Where(&model.User{UserName: u.UserName}).First(&user).Error
	if err == nil {
		return errors.New("用户已存在")
	} else {
		//处理角色
		var roleArr []string
		for _, v := range u.RoleGroup {
			roleArr = append(roleArr, v.RoleName)
		}
		roles, err := FindRoleIdsByRoleNameArr(roleArr)
		if err != nil {
			return err
		}
		u.RoleGroup = roles
		u.UUID = uuid.NewV4()
		u.Password = encrypt_plugin.BcryptEncode(u.Password)
		u.InvitationCode = encrypt_plugin.RandomString(8)
		return CreateUser(NewUserSubscribe(&u))
	}
}

// 新注册用户分配套餐
func NewUserSubscribe(u *model.User) *model.User {
	//查询商品信息
	if global.Server.System.DefaultGoods == "" {
		return u
	}
	var goods = model.Goods{
		Subject: global.Server.System.DefaultGoods,
	}
	//查询默认套餐
	g, _, err := CommonSqlFind[model.Goods, model.Goods, model.Goods](goods)
	if err != nil {
		return u
	}
	// 处理用户订阅信息
	return HandleUserSubscribe(u, &g)
}

// 用户登录
func Login(u *model.UserLogin) (*model.User, error) {
	var user model.User
	err := global.DB.Where("user_name = ?", u.UserName).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("用户不存在")
	} else if !user.Enable {
		return nil, errors.New("用户已冻结")
	} else {
		if err := encrypt_plugin.BcryptDecode(u.Password, user.Password); err != nil {
			return nil, errors.New("密码错误")
		}
		return &user, err
	}
	return &user, err
}

// 获取当前请求节点（根据 node_id 参数判断）可连接的用户
func FindUsersByGoods(goods *[]model.Goods) (*[]model.SSUsers, error) {
	var goodsArr []int64
	for _, v := range *goods {
		goodsArr = append(goodsArr, v.ID)
	}
	var users []model.SSUsers
	//err := global.DB.Model(&model.User{}).Where("goods_id in (?)", goodsArr).Find(&users).Error
	err := global.DB.Model(&model.User{}).Where("goods_id in (?) and sub_status = ?", goodsArr, true).Find(&users).Error
	return &users, err
}

// 查询订单属于哪个用户
func FindUsersByOrderID(outTradeNo string) (*model.User, error) {
	var order model.Orders
	err := global.DB.Where("out_trade_no = ?", outTradeNo).Preload("User").Find(&order).Error
	return &order.User, err
}

// 查用户 by user_id
func FindUserByID(id int64) (*model.User, error) {
	var u model.User
	err := global.DB.First(&u, id).Error
	return &u, err
}

// 更新用户订阅信息
func UpdateUserSubscribe(order *model.Orders) error {
	//查询商品信息
	goods, _ := FindGoodsByGoodsID(order.GoodsID)
	//查询订单属于哪个用户
	u, _ := FindUserByID(order.UserID)
	//构建用户订阅信息
	user := HandleUserSubscribe(u, goods)
	//更新用户订阅信息
	return global.DB.Save(&user).Error
}

// 处理用户订阅信息
func HandleUserSubscribe(u *model.User, goods *model.Goods) *model.User {
	u.SubscribeInfo.T = goods.TotalBandwidth * 1024 * 1024 * 1024 // TotalBandwidth单位：GB。总流量单位：B
	u.SubscribeInfo.U = 0
	u.SubscribeInfo.D = 0
	if u.SubscribeInfo.SubscribeUrl == "" {
		u.SubscribeInfo.SubscribeUrl = encrypt_plugin.RandomString(8) //随机字符串订阅url
	}
	u.SubscribeInfo.GoodsID = goods.ID           //当前订购的套餐
	u.SubscribeInfo.GoodsSubject = goods.Subject //套餐标题
	u.SubscribeInfo.SubStatus = true             //订阅状态
	t := time.Now().AddDate(0, 0, int(goods.ExpirationDate))
	u.SubscribeInfo.ExpiredAt = &t //过期时间
	if goods.NodeConnector != 0 {
		u.SubscribeInfo.NodeConnector = goods.NodeConnector //连接客户端数
	}
	return u
}

// 批量更新用户流量信息
func UpdateUserTrafficInfo(userArr []model.User, userIds []int64) error {
	var userArrQuery []model.User
	err := global.DB.Where("id in ?", userIds).Select("id", "u", "d").Find(&userArrQuery).Error
	if err != nil {
		return err
	}
	for item, _ := range userArrQuery {
		userArrQuery[item].SubscribeInfo.U = userArrQuery[item].SubscribeInfo.U + userArr[item].SubscribeInfo.U
		userArrQuery[item].SubscribeInfo.D = userArrQuery[item].SubscribeInfo.D + userArr[item].SubscribeInfo.D
	}
	return global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"u", "d"}),
	}).Create(&userArrQuery).Error

}

// 用户流量，有效期 检测任务
func UserExpiryCheck() error {
	//fmt.Println("开始用户流量，有效期 检测任务")
	return global.DB.Exec("update user set sub_status = 0 where expired_at < ? or ( u + d ) > t", time.Now()).Error
}

// 修改混淆
func ChangeSubHost(uID int64, host string) error {
	u := map[string]any{
		"host": host,
	}
	return global.DB.Model(&model.User{ID: uID}).Updates(u).Error
}

// 获取自身信息
func GetUserInfo(uID int64) (*model.User, error) {
	var user model.User
	return &user, global.DB.First(&user, uID).Error
}

// 获取用户列表,分页
func GetUserlist(params *model.PaginationParams) (*model.UsersWithTotal, error) {
	var userArr model.UsersWithTotal
	var err error
	if params.Search != "" {
		err = global.DB.Model(&model.User{}).Where("user_name like ?", ("%" + params.Search + "%")).Count(&userArr.Total).Limit(int(params.PageSize)).Offset((int(params.PageNum) - 1) * int(params.PageSize)).Preload("RoleGroup").Find(&userArr.UserList).Error
	} else {
		err = global.DB.Model(&model.User{}).Count(&userArr.Total).Limit(int(params.PageSize)).Offset((int(params.PageNum) - 1) * int(params.PageSize)).Preload("RoleGroup").Find(&userArr.UserList).Error
	}
	return &userArr, err
}

// 更新用户信息
func UpdateUser(u *model.User) error {
	// $2a$10$
	if !strings.HasPrefix(u.Password, "$2a$10$") {
		u.Password = encrypt_plugin.BcryptEncode(u.Password)
	}
	return global.DB.Updates(&u).Error
}
func SaveUser(u *model.User) error {
	// $2a$10$
	if !strings.HasPrefix(u.Password, "$2a$10$") {
		u.Password = encrypt_plugin.BcryptEncode(u.Password)
	}
	return global.DB.Save(&u).Error
}

// 创建用户
func CreateUser(u *model.User) error {
	return global.DB.Create(&u).Error
}

// 删除用户
func DeleteUser(u *model.User) error {
	return global.DB.Delete(&u).Error
}

// 重置用户密码
func ResetUserPassword(u *model.User) error {
	return global.DB.Model(&model.User{}).Where("user_name = ?", u.UserName).Updates(&u).Error
}

// 重置管理员密码
func ResetAdminPassword() {
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
		SaveUser(&user)
		fmt.Println("完成...")

	case "n":
		os.Exit(0)
	}

}

// 处理推荐人返利
func ReferrerRebate(uID int64, receiptAmount string) {
	u, err := FindUserByID(uID)
	if err != nil || len(u.ReferrerCode) < 8 { //error或者推荐人为空
		return
	}
	var referrerUser model.User //查找推荐人信息
	err = global.DB.Where(&model.User{InvitationCode: u.ReferrerCode}).First(&referrerUser).Error
	if err != nil {
		return //推荐人不存在
	}
	a, _ := strconv.ParseFloat(receiptAmount, 64)
	referrerUser.Remain = referrerUser.Remain + a*global.Server.System.RebateRate
	SaveUser(&referrerUser)
}

// 处理用户余额
func RemainHandle(uid int64, remain string) {
	remainFloat64, _ := strconv.ParseFloat(remain, 64)
	if remainFloat64 == 0 {
		return
	}
	user, _ := FindUserByID(uid)
	user.Remain = user.Remain - remainFloat64
	SaveUser(user)
}
