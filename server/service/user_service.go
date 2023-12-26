package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"strings"
	"time"

	"errors"
	"github.com/ppoonk/AirGo/model"
	encrypt_plugin "github.com/ppoonk/AirGo/utils/encrypt_plugin"
	uuid "github.com/satori/go.uuid"
)

// 注册
func Register(u *model.User) error {
	//判断是否存在
	var user model.User
	err := global.DB.Where(&model.User{UserName: u.UserName}).First(&user).Error
	if err == nil {
		return errors.New("User already exists")
	} else if err == gorm.ErrRecordNotFound {
		var newUser = model.User{
			UUID:           uuid.NewV4(),
			UserName:       u.UserName,
			NickName:       u.UserName,
			Avatar:         u.Avatar,                                //头像
			Password:       encrypt_plugin.BcryptEncode(u.Password), //密码
			RoleGroup:      []model.Role{{ID: 2}},                   //默认角色：普通用户角色
			InvitationCode: encrypt_plugin.RandomString(8),          //邀请码
			ReferrerCode:   u.ReferrerCode,                          //推荐人
			SubscribeInfo: model.SubscribeInfo{
				SubscribeUrl: encrypt_plugin.RandomString(8), //随机字符串订阅url
			},
		}
		//通知
		if global.Server.Notice.WhenUserRegistered {
			global.GoroutinePool.Submit(func() {
				UnifiedPushMessage("新注册用户：" + newUser.UserName)
			})
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
		return errors.New("User already exists")
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
	if global.Server.Subscribe.DefaultGoods == 0 {
		return u
	}
	var goods = model.Goods{
		ID: global.Server.Subscribe.DefaultGoods,
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
		return nil, errors.New("User does not exist")
	} else if !user.Enable {
		return nil, errors.New("User frozen")
	} else {
		if err := encrypt_plugin.BcryptDecode(u.Password, user.Password); err != nil {
			return nil, errors.New("Password error")
		}
		return &user, err
	}
	return &user, err
}

// 查用户 by user_id
func FindUserByID(id int64) (*model.User, error) {
	var u model.User
	err := global.DB.First(&u, id).Error
	return &u, err
}

// 查用户 by tg_id
func FindUserByTgID(tgID int64) (*model.User, error) {
	var u model.User
	err := global.DB.Where("tg_id = ?", tgID).First(&u).Error
	return &u, err
}
func FindUserByUserName(userName string) (*model.User, error) {
	var u model.User
	err := global.DB.Where("user_name = ?", userName).First(&u).Error
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
	u.SubscribeInfo.GoodsID = goods.ID           //当前订购的套餐
	u.SubscribeInfo.GoodsSubject = goods.Subject //套餐标题
	u.SubscribeInfo.SubStatus = true             //订阅状态
	if u.SubscribeInfo.SubscribeUrl == "" {
		u.SubscribeInfo.SubscribeUrl = encrypt_plugin.RandomString(8) //随机字符串订阅url
	}
	if goods.NodeConnector != 0 {
		u.SubscribeInfo.NodeConnector = goods.NodeConnector //连接客户端数
	}
	switch goods.TrafficResetMethod { //判断是否叠加
	case "Stack":
		u.SubscribeInfo.T = u.SubscribeInfo.T + goods.TotalBandwidth*1024*1024*1024 // GB->MB->KB->B
		t := u.SubscribeInfo.ExpiredAt.AddDate(0, 0, int(goods.ExpirationDate))
		u.SubscribeInfo.ExpiredAt = &t //过期时间
	default:
		u.SubscribeInfo.T = goods.TotalBandwidth * 1024 * 1024 * 1024 // GB->MB->KB->B
		u.SubscribeInfo.U = 0
		u.SubscribeInfo.D = 0
		t := time.Now().AddDate(0, 0, int(goods.ExpirationDate))
		u.SubscribeInfo.ExpiredAt = &t //过期时间
	}
	u.SubscribeInfo.ResetDay = goods.ResetDay //流量重置日
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
func UpdateUserTrafficLog(UserTrafficLogMap map[int64]model.UserTrafficLog, userIds []int64) error {
	var query []model.UserTrafficLog
	now := time.Now()
	//当日0点
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := global.DB.Where("created_at > ? AND user_id IN ?", todayZero, userIds).Find(&query).Error
	if err != nil {
		return err
	}
	for k, _ := range query {
		if tl, ok := UserTrafficLogMap[query[k].UserID]; ok { //已存在，叠加流量
			query[k].U += tl.U
			query[k].D += tl.D
			delete(UserTrafficLogMap, query[k].UserID) //删除
		}
	}
	//不存在的数据，追加到最后面，一起插入数据库
	if len(UserTrafficLogMap) > 0 {
		for k, _ := range UserTrafficLogMap {
			query = append(query, UserTrafficLogMap[k])
		}
	}
	if len(query) == 0 {
		return nil
	}
	return global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"u", "d"}),
	}).Create(&query).Error

}

// 用户流量，有效期 检测任务
func UserExpiryCheck() error {
	err := global.DB.Exec("UPDATE user SET sub_status = 0 WHERE expired_at < ? or ( u + d ) > t", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

// 用户流量重置任务
func UserTrafficReset() error {
	day := time.Now().Day()
	err := global.DB.Exec("UPDATE user SET u = 0, d = 0 WHERE reset_day = ?", day).Error
	if err != nil {
		return err
	}
	return nil
}

// 修改混淆
func ChangeSubHost(uID int64, host string) error {
	u := map[string]any{
		"host": host,
	}
	return global.DB.Model(&model.User{ID: uID}).Updates(u).Error
}

// 获取用户列表
func GetUserlist(params *model.FieldParamsReq) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var userList []model.User
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	err := global.DB.Model(&model.User{}).Count(&data.Total).Where(dataSql).Preload("RoleGroup").Find(&userList).Error
	if err != nil {
		global.Logrus.Error("GetUserlist error:", err.Error())
		return nil, err
	}
	data.Data = userList
	return &data, nil
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
	//删除关联的角色组
	err := DeleteUserRoleGroup(u)
	if err != nil {
		return err
	}
	//删除关联的订单
	err = DeleteUserAllOrder(u)
	if err != nil {
		return err
	}
	//删除用户
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
	referrerUser.Remain = referrerUser.Remain + a*global.Server.Subscribe.RebateRate
	SaveUser(&referrerUser)
}

// 处理用户余额
func RemainHandle(uid int64, remain string) error {
	remainFloat64, _ := strconv.ParseFloat(remain, 64)
	if remainFloat64 == 0 {
		return nil
	}
	user, err := FindUserByID(uid)
	if err != nil {
		return err
	}
	user.Remain = user.Remain - remainFloat64
	return SaveUser(user)
}

// 处理用户充值卡商品
func RechargeHandle(order *model.Orders) error {
	Show(order)
	//查询商品信息
	goods, _ := FindGoodsByGoodsID(order.GoodsID)
	Show(goods)
	orderRemainAmount, _ := strconv.ParseFloat(order.RemainAmount, 64)
	rechargeFloat64, _ := strconv.ParseFloat(goods.RechargeAmount, 64)
	user, err := FindUserByID(order.UserID)
	if err != nil {
		return err
	}
	user.Remain = user.Remain - orderRemainAmount + rechargeFloat64
	if user.Remain < 0 {
		user.Remain = 0
	}
	Show(order)
	return SaveUser(user)
}

// 打卡
func ClockIn(uID int64) (int, int, error) {
	//查询用户信息
	user, err := FindUserByID(uID)
	if err != nil {
		return 0, 0, err
	}
	//判断订阅是否有效
	if !user.SubscribeInfo.SubStatus {
		return 0, 0, errors.New("subscribe is expired")
	}
	//随机流量
	t := encrypt_plugin.RandomNumber(int(global.Server.Subscribe.ClockInMinTraffic), int(global.Server.Subscribe.ClockInMaxTraffic)) //MB
	user.SubscribeInfo.T = int64(t)*1024*1024 + user.SubscribeInfo.T
	//随机天数
	day := encrypt_plugin.RandomNumber(int(global.Server.Subscribe.ClockInMinDay), int(global.Server.Subscribe.ClockInMaxDay))
	*user.SubscribeInfo.ExpiredAt = user.SubscribeInfo.ExpiredAt.AddDate(0, 0, day)

	err = SaveUser(user)
	return t, day, err

}

func GetUserTraffic(params *model.FieldParamsReq) (*model.UserTrafficLog, error) {
	var userTraffic model.UserTrafficLog
	var err error
	_, dataSql := CommonSqlFindNoOrderByNoLimitSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		return nil, errors.New("invalid parameter")
	}
	//fmt.Println("dataSql:", dataSql)
	if global.Config.SystemParams.DbType == "mysql" {
		err = global.DB.Model(&model.UserTrafficLog{}).Where(dataSql).Select("user_id, any_value(user_name) AS user_name, SUM(u) AS u, SUM(d) AS d").Group("user_id").Find(&userTraffic).Error
	} else {
		err = global.DB.Model(&model.UserTrafficLog{}).Where(dataSql).Select("user_id, user_name, SUM(u) AS u, SUM(d) AS d").Group("user_id").Find(&userTraffic).Error
	}

	return &userTraffic, err
}

func GetAllUserTraffic(params *model.FieldParamsReq) (*model.CommonDataResp, error) {
	//约定：params.FieldParamsList 数组前两项传时间，第三个开始传查询参数
	var userTraffic []model.UserTrafficLog
	var total int64
	var err error
	_, dataSql := CommonSqlFindNoOrderByNoLimitSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		return nil, errors.New("invalid parameter")
	}
	//fmt.Println("dataSql:", dataSql)

	if global.Config.SystemParams.DbType == "mysql" { //mysql only_full_group_by 问题
		err = global.DB.Model(&model.UserTrafficLog{}).Where(dataSql).Select("user_id, any_value(user_name) AS user_name, SUM(u) u, SUM(d) AS d").Group("user_id").Count(&total).Order(params.Pagination.OrderBy).Limit(int(params.Pagination.PageSize)).Offset((int(params.Pagination.PageNum) - 1) * int(params.Pagination.PageSize)).Find(&userTraffic).Error
	} else {
		err = global.DB.Model(&model.UserTrafficLog{}).Where(dataSql).Select("user_id, user_name, SUM(u) u, SUM(d) AS d").Group("user_id").Count(&total).Order(params.Pagination.OrderBy).Limit(int(params.Pagination.PageSize)).Offset((int(params.Pagination.PageNum) - 1) * int(params.Pagination.PageSize)).Find(&userTraffic).Error
	}

	if err != nil {
		return nil, err
	}
	return &model.CommonDataResp{
		Total: total,
		Data:  userTraffic,
	}, nil
}

// 临时代码，删除用户流量统计
func DeleteUserTrafficTemp() error {
	return global.DB.Where("id > 0").Delete(&model.UserTrafficLog{}).Error

}
