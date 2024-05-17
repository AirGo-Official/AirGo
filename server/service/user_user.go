package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/jwt_plugin"
	"github.com/ppoonk/AirGo/utils/time_plugin"
	timeTool "github.com/ppoonk/AirGo/utils/time_plugin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"errors"
	"github.com/ppoonk/AirGo/model"
	encrypt_plugin "github.com/ppoonk/AirGo/utils/encrypt_plugin"
)

type User struct{}

var UserSvc *User

// 注册
func (us *User) Register(userParams *model.User) error {
	//判断是否存在
	var user model.User
	err := global.DB.Where(&model.User{UserName: userParams.UserName}).First(&user).Error
	if err == nil {
		return errors.New("User already exists")
	} else if err == gorm.ErrRecordNotFound {

		return us.CreateUser(userParams)
	} else {
		return err
	}
}

// 创建用户
func (us *User) CreateUser(u *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&u).Error
	})
}

// 用户登录
func (us *User) Login(u *model.UserLoginRequest) (*model.User, error) {
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
}
func (us *User) GetUserToken(user *model.User) (string, error) {
	//查缓存
	cache, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, user.ID))
	if ok {
		return cache.(string), nil
	}
	//生成新的
	myCustomClaimsPrefix := jwt_plugin.MyCustomClaimsPrefix{
		UserID:   user.ID,
		UserName: user.UserName,
	}
	ep, _ := timeTool.ParseDuration(global.Server.Security.JWT.ExpiresTime)
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    global.Server.Security.JWT.Issuer,      // 签发者
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)), //过期时间
		NotBefore: jwt.NewNumericDate(time.Now()),         //生效时间
	}
	tokenNew, err := jwt_plugin.GenerateTokenUsingHs256(myCustomClaimsPrefix, registeredClaims, global.Server.Security.JWT.SigningKey)
	if err != nil {
		return "", err
	}
	global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, user.ID), tokenNew, ep)
	return tokenNew, nil
}

// 查用户
func (us *User) FirstUser(user *model.User) (*model.User, error) {
	var userQuery model.User
	err := global.DB.Where(&user).First(&userQuery).Error
	return &userQuery, err
}

// 保存用户信息
func (us *User) SaveUser(u *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&u).Error
	})
}

// 更新用户信息
func (us *User) UpdateUser(userParams *model.User, values map[string]any) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Debug().Model(&model.User{}).Where(&userParams).Updates(values).Error
	})
}

// 处理用户充值卡商品
func (us *User) RechargeHandle(order *model.Order) error {
	//查询商品信息
	goods, _ := ShopSvc.FirstGoods(&model.Goods{ID: order.GoodsID})
	rechargeFloat64, _ := strconv.ParseFloat(goods.RechargeAmount, 64)
	user, err := us.FirstUser(&model.User{ID: order.UserID})
	if err != nil {
		return err
	}
	startAmount := user.Balance
	endAmount := fmt.Sprintf("%.2f", user.Balance+rechargeFloat64)
	endBalance, _ := strconv.ParseFloat(endAmount, 64)
	user.Balance = endBalance
	if user.Balance < 0 {
		user.Balance = 0
	}
	err = us.SaveUser(user)
	if err != nil {
		return err
	}
	//处理余额流水
	err = FinanceSvc.NewBalanceStatement(&model.BalanceStatement{
		UserID:      order.UserID,
		Title:       constant.BALANCE_STATEMENT_TITLE_RECHARGE,
		Type:        constant.BALANCE_STATEMENT_TYPE_PLUS,
		Amount:      goods.RechargeAmount,
		FinalAmount: endAmount,
	})
	if err != nil {
		return err
	}
	if user.WhenBalanceChanged { //通知
		global.GoroutinePool.Submit(func() {
			us.PushMessageWhenBalanceChanged(user, startAmount, rechargeFloat64)
		})
	}
	return nil
}

// 处理余额支付
func (us *User) UserBalancePayHandler(order *model.Order) error {
	// 查询user
	user, err := us.FirstUser(&model.User{ID: order.UserID})
	if err != nil {
		return err
	}
	startAmount := user.Balance
	totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
	if totalAmount == 0 {
		return nil
	}
	if user.Balance < totalAmount {
		return errors.New(constant.ERROR_BALANCE_IS_NOT_ENOUGH)
	}
	endAmount := fmt.Sprintf("%.2f", user.Balance-totalAmount)
	res, _ := strconv.ParseFloat(endAmount, 64)
	user.Balance = res
	err = us.SaveUser(user)
	if err != nil {
		return err
	}
	//处理余额流水
	err = FinanceSvc.NewBalanceStatement(&model.BalanceStatement{
		UserID:      order.UserID,
		Title:       constant.BALANCE_STATEMENT_TITLE_EXPENDITURE,
		Type:        constant.BALANCE_STATEMENT_TYPE_REDUCE,
		Amount:      order.TotalAmount,
		FinalAmount: endAmount,
	})
	if err != nil {
		return err
	}
	if user.WhenBalanceChanged { //通知
		global.GoroutinePool.Submit(func() {
			us.PushMessageWhenBalanceChanged(user, startAmount, totalAmount)
		})
	}
	return nil
}

// 通知
func (us *User) PushMessageWhenBalanceChanged(user *model.User, startAmount, changedAmount float64) {
	msg := MessageInfo{
		UserID:      user.ID,
		MessageType: MESSAGE_TYPE_USER,
		User:        user,
		Message: strings.Join([]string{
			"【余额变动提醒】",
			fmt.Sprintf("时间：%s", time.Now().Format("2006-01-02 15:04:05")),
			fmt.Sprintf("开始余额：%s", fmt.Sprintf("%.2f", startAmount)),
			fmt.Sprintf("结束余额：%s", fmt.Sprintf("%.2f", user.Balance)),
			fmt.Sprintf("变动值：%s\n", fmt.Sprintf("%.2f", changedAmount)),
		}, "\n"),
	}
	PushMessageSvc.PushMessage(&msg)
}

// 删除cache
func (us *User) DeleteUserCacheTokenByID(user *model.User) {
	global.LocalCache.Delete(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, user.ID))
}

// 校验注册时的邮箱验证码
func (us *User) VerifyEmailWhenRegister(params model.UserRegister) (bool, error) {
	//处理邮箱验证码
	userEmail := params.UserName + params.EmailSuffix //处理邮箱后缀,注册时，用户名和邮箱后缀是分开的
	cacheEmail, ok := global.LocalCache.Get(constant.CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME + userEmail)
	if ok {
		if !strings.EqualFold(cacheEmail.(string), params.EmailCode) {
			//验证失败，返回错误响应，但不删除缓存的验证码。因为用户输错了，需要重新输入，而不需要重新发送验证码
			return false, nil
		} else {
			//验证成功，删除缓存的验证码
			global.LocalCache.Delete(constant.CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME + userEmail)
			return true, nil
		}
	} else {
		//cache缓存超时
		return false, errors.New("The verification code has expired, please try again")
	}

}

// 校验重置密码时的邮箱验证码
func (us *User) VerifyEmailWhenResetPassword(params model.UserLoginRequest) (bool, error) {
	cacheEmail, ok := global.LocalCache.Get(constant.CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME + params.UserName)
	if ok {
		if !strings.EqualFold(cacheEmail.(string), params.EmailCode) {
			//验证失败，返回错误响应，但不删除缓存的验证码。因为用户输错了，需要重新输入，而不需要重新发送验证码
			return false, nil
		} else {
			//验证成功，删除缓存的验证码
			global.LocalCache.Delete(constant.CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME + params.UserName)
			return true, nil
		}
	} else {
		//cache缓存超时
		return false, errors.New("The verification code has expired, please try again")
	}
}

// 打卡抽奖
func (us *User) ClockIn(uID int64) (int, float64, error) {
	if !global.Server.Finance.EnableLottery {
		return 0, 0, errors.New(constant.ERROR_SERVICE_NOT_ENABLED)
	}
	//判断是否已签到打卡
	_, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_USER_IS_CLOCKIN_BY_ID, uID))
	if ok {
		return 0, 0, errors.New("You have already drawn a lottery today")
	}
	var (
		totalWeight int
		index       int     //中奖索引
		balance     float64 //奖励余额
	)
	for _, v := range global.Server.Finance.Jackpot {
		totalWeight += v.Weight
	}
	num := encrypt_plugin.RandomNumber(0, int(totalWeight))
	for k, v := range global.Server.Finance.Jackpot {
		if num < v.Weight {
			index = k
			balance = v.Balance
			break
		}
		num -= v.Weight
	}
	user, err := us.FirstUser(&model.User{ID: uID})
	if err != nil {
		return 0, 0, err
	}
	startAmount := user.Balance
	user.Balance += balance
	//格式化2位小数
	user.Balance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", user.Balance), 64)
	err = us.SaveUser(user)
	if err != nil {
		return 0, 0, err
	}
	//处理余额流水
	err = FinanceSvc.NewBalanceStatement(&model.BalanceStatement{
		UserID:      uID,
		Title:       constant.BALANCE_STATEMENT_TITLE_PRIZE, //类型为奖励
		Type:        constant.BALANCE_STATEMENT_TYPE_PLUS,
		Amount:      fmt.Sprintf("%.2f", balance),
		FinalAmount: fmt.Sprintf("%.2f", user.Balance),
	})
	if err != nil {
		return 0, 0, err
	}
	if user.WhenBalanceChanged { //通知
		global.GoroutinePool.Submit(func() {
			us.PushMessageWhenBalanceChanged(user, startAmount, user.Balance)
		})
	}
	//用户今天已经打卡
	global.LocalCache.Set(fmt.Sprintf("%s%d",
		constant.CACHE_USER_IS_CLOCKIN_BY_ID, uID),
		struct{}{},
		time_plugin.GetTimeIntervalBetweenNowAndMidnightTheNextDay())

	return index, balance, nil
}
