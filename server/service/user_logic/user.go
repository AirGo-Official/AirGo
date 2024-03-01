package user_logic

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/jwt_plugin"
	timeTool "github.com/ppoonk/AirGo/utils/time_plugin"
	"gorm.io/gorm"
	"strconv"
	"time"

	"errors"
	"github.com/ppoonk/AirGo/model"
	encrypt_plugin "github.com/ppoonk/AirGo/utils/encrypt_plugin"
)

type User struct{}

var userService *User

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
	return global.DB.Create(&u).Error
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

// 更新用户信息
func (us *User) UpdateUser(userParams *model.User, values map[string]any) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.User{}).Where(&userParams).Updates(values).Error
	})
}

// 处理用户充值卡商品
func (us *User) RechargeHandle(order *model.Order) error {
	//查询商品信息
	goods, _ := shopService.FirstGoods(&model.Goods{ID: order.GoodsID})
	orderRemainAmount, _ := strconv.ParseFloat(order.BalanceAmount, 64)
	rechargeFloat64, _ := strconv.ParseFloat(goods.RechargeAmount, 64)
	user, err := us.FirstUser(&model.User{ID: order.UserID})
	if err != nil {
		return err
	}
	user.Balance = user.Balance - orderRemainAmount + rechargeFloat64
	if user.Balance < 0 {
		user.Balance = 0
	}
	return us.SaveUser(user)
}

// 处理余额支付
func (us *User) UserRemainPayHandler(order *model.Order) error {
	// 查询user
	user, err := us.FirstUser(&model.User{ID: order.UserID})
	if err != nil {
		return err
	}
	totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
	if user.Balance < totalAmount {
		return errors.New(constant.ERROR_BALANCE_IS_NOT_ENOUGH)
	}
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", user.Balance-totalAmount), 64)
	user.Balance = res
	return userService.SaveUser(user)
}

// 保存用户信息
func (us *User) SaveUser(u *model.User) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&u).Error
	})
}

// 删除cache
func (us *User) DeleteUserCacheTokenByID(user *model.User) {
	global.LocalCache.Delete(fmt.Sprintf("%s%d", constant.CACHE_USER_TOKEN_BY_ID, user.ID))
}
