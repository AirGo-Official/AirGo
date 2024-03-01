package user_logic

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Order struct {
}

var orderService *Order

func InitOrderSvc() {
	orderService.StartTask()
}

func (o *Order) StartTask() {
	ch, err := global.Queue.Subscribe(constant.SUBMIT_ORDER, 1000) //1000 缓冲
	if err != nil {
		global.Logrus.Error(err)
		return
	}
	go func() {
		for v := range ch {
			preOrder := v.(*model.Order)
			switch preOrder.OrderType {
			case constant.ORDER_TYPE_NEW:
				// 1、查找商品
				goods, err := shopService.FirstGoods(&model.Goods{ID: preOrder.GoodsID})
				if err != nil {
					continue
				}
				// 2、检查库存
				if goods.Stock <= 0 {
					continue
				}
				// 3、扣库存
				goods.Stock--
				// 4、更新商品
				err = shopService.UpdateGoods(goods)
				if err != nil {
					continue
				}
				// 5、生成订单,订单状态：待支付
				preOrder.CreatedAt = time.Now()
				preOrder.TradeStatus = constant.ORDER_STATUS_WAIT_BUYER_PAY
				// 6、存入数据库、cache
				err = o.CreateOrder(preOrder)
				if err != nil {
					continue
				}
				// 7、设置5分钟过期，5分钟等待付款
				global.LocalCache.
					Set(constant.CACHE_SUBMIT_ORDER_BY_ORDERID+preOrder.OutTradeNo,
						preOrder,
						constant.CACHE_SUBMIT_ORDER_TIMEOUT*time.Minute)
				// 8、删除幂等性cache
				global.LocalCache.
					Delete(fmt.Sprintf("%s%d:%d", constant.CACHE_USERID_AND_GOODSID, preOrder.UserID, preOrder.GoodsID))

			case constant.ORDER_TYPE_RENEW:
				// 1、生成订单,订单状态：待支付
				preOrder.CreatedAt = time.Now()
				preOrder.TradeStatus = constant.ORDER_STATUS_WAIT_BUYER_PAY
				// 2、存入数据库、cache
				err = o.CreateOrder(preOrder)
				if err != nil {
					continue
				}
				// 3、设置5分钟过期，5分钟等待付款
				global.LocalCache.
					Set(constant.CACHE_SUBMIT_ORDER_BY_ORDERID+preOrder.OutTradeNo,
						preOrder,
						constant.CACHE_SUBMIT_ORDER_TIMEOUT*time.Minute)
				// 4、删除幂等性cache
				global.LocalCache.
					Delete(fmt.Sprintf("%s%d:%d", constant.CACHE_USERID_AND_CUSTOMERSERVICEID, preOrder.UserID, preOrder.CustomerServiceID))

			case constant.ORDER_TYPE_RESTORE:
				// 判断是否是续费订单
				if preOrder.CustomerServiceID == 0 {
					goods, err := shopService.FirstGoods(&model.Goods{ID: preOrder.GoodsID})
					if err == nil {
						// 库存补偿
						goods.Stock++
						// 更新商品
						_ = shopService.UpdateGoods(goods)
					}
				}
				// 更新订单，订单关闭
				preOrder.TradeStatus = constant.ORDER_STATUS_TRADE_CLOSED
				_ = o.UpdateOrder(preOrder)
			}
		}
	}()

}

func (o *Order) PreCheckOrder(orderReq *model.Order) error {
	switch orderReq.OrderType {
	case constant.ORDER_TYPE_NEW:
		// 幂等性
		_, ok := global.LocalCache.Get(fmt.Sprintf("%s%d:%d", constant.CACHE_USERID_AND_GOODSID, orderReq.UserID, orderReq.GoodsID))
		if ok {
			return errors.New(constant.ERROR_DUPLICATE_ORDER)
		}
		// 库存
		goods, err := shopService.FirstGoods(&model.Goods{ID: orderReq.GoodsID})
		if err != nil {
			return err
		}
		if goods.Stock <= 0 {
			return errors.New(constant.ERROR_STOCK_OF_GOODS_EMPTY)
		}
		// 限购
		err, ok = o.CheckQuota(orderReq, goods)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New(constant.ERROR_GOODS_EXCEEDED_QUOTA)
		}
		return nil

	case constant.ORDER_TYPE_RENEW:
		// 幂等性
		_, ok := global.LocalCache.Get(fmt.Sprintf("%s%d:%d", constant.CACHE_USERID_AND_CUSTOMERSERVICEID, orderReq.UserID, orderReq.CustomerServiceID))
		if ok {
			return errors.New(constant.ERROR_DUPLICATE_ORDER)
		}
		// 检查服务
		cs, err := customerService.FirstCustomerService(&model.CustomerService{ID: orderReq.CustomerServiceID, UserID: orderReq.UserID})
		if err != nil {
			return err
		}
		//是否允许续费
		if !cs.IsRenew {
			return errors.New(constant.ERROR_CUSTOMER_SERVICE_NO_RENEWAL)
		}
		// 检查关联的商品
		goods, err := shopService.FirstGoods(&model.Goods{ID: orderReq.GoodsID})
		if err != nil {
			return err
		}
		//原商品是否在售
		if !goods.IsSale {
			return errors.New(constant.ERROR_GOODS_NOT_SALE)
		}
		return nil
	default:
		return errors.New(constant.ERROR_INVALID_ORDER_TYPE)
	}
}
func (o *Order) CheckQuota(orderReq *model.Order, goods *model.Goods) (error, bool) {
	var total int64
	var orders []model.Order
	err := global.DB.Model(&model.Order{}).Where(&model.Order{OrderType: constant.ORDER_TYPE_NEW, TradeStatus: constant.ORDER_STATUS_TRADE_SUCCESS, UserID: orderReq.UserID, GoodsID: goods.ID}).
		Count(&total).
		Find(&orders).Error
	if err != nil {
		return err, false
	}
	//fmt.Println("total:", total, "goods.Quota:", goods.Quota)
	return nil, total < goods.Quota
}

// 订单预处理，计算价格
func (o *Order) PreHandleOrder(orderReq *model.Order) (*model.Order, string, error) {
	// 判断订单类型
	var preOrder model.Order
	var msg string
	switch orderReq.OrderType {
	case constant.ORDER_TYPE_NEW:
		//通过商品id查找商品
		goods, err := shopService.FirstGoods(&model.Goods{ID: orderReq.GoodsID})
		if err != nil {
			return nil, "", err
		}
		//构造系统订单参数
		price, err := strconv.ParseFloat(goods.Price, 64)
		if err != nil {
			return nil, "", err
		}
		if orderReq.Duration <= 0 {
			orderReq.Duration = 1
		}
		originalAmount := fmt.Sprintf("%.2f", price*float64(orderReq.Duration))
		preOrder = model.Order{
			OrderType:      constant.ORDER_TYPE_NEW,
			TradeStatus:    constant.ORDER_STATUS_CREATED,
			OutTradeNo:     time.Now().Format("20060102150405") + fmt.Sprintf("%d", orderReq.UserID), //系统订单号：时间戳+user id]
			OriginalAmount: originalAmount,                                                           //原始价格
			TotalAmount:    originalAmount,                                                           //订单价格
			BuyerPayAmount: "0.00",
			CouponAmount:   "0.00",
			BalanceAmount:  "0.00",

			UserID:   orderReq.UserID,
			UserName: orderReq.UserName,
			//User:           model.User{},

			GoodsID:     goods.ID,
			GoodsType:   goods.GoodsType,
			DeliverType: goods.DeliverType,
			//DeliverText: "",
			Subject:  goods.Subject,
			Price:    goods.Price,
			Duration: orderReq.Duration,

			//CustomerServiceID: 0,

			//PayID:           0,
			//PayType:         "",
			//PayInfo:         model.PreCreatePayToFrontend{},
			//TradeNo:        "",
			//BuyerLogonId:    "",
			//CouponID:        receiveOrder.CouponID,
			CouponName: orderReq.CouponName,
		}
	case constant.ORDER_TYPE_RENEW:
		// 查找用户服务
		cs, err := customerService.FirstCustomerService(&model.CustomerService{UserID: orderReq.UserID, ID: orderReq.CustomerServiceID})
		if err != nil {
			return nil, "", err
		}
		//fmt.Println("构造系统订单参数:")
		//Show(cs)
		// 构造系统订单参数
		preOrder = model.Order{
			OrderType:      constant.ORDER_TYPE_RENEW,
			TradeStatus:    constant.ORDER_STATUS_CREATED,
			OutTradeNo:     time.Now().Format("20060102150405") + fmt.Sprintf("%d", orderReq.UserID), //系统订单号：时间戳+user id,
			OriginalAmount: cs.RenewalAmount,
			TotalAmount:    cs.RenewalAmount,
			BuyerPayAmount: "0.00",
			CouponAmount:   "0.00",
			BalanceAmount:  "0.00",

			UserID:   orderReq.UserID,
			UserName: orderReq.UserName,
			//User:           model.User{},

			GoodsID:   cs.GoodsID,
			GoodsType: cs.GoodsType,
			//DeliverType: "",
			//DeliverText: "",
			Subject:  cs.Subject,
			Price:    cs.Price,
			Duration: cs.Duration,

			CustomerServiceID: cs.ID,

			//PayID:        0,
			//PayType:      "",
			//PayInfo:      model.PreCreatePayToFrontend{},
			//TradeNo:      "",
			//BuyerLogonId: "",
			//CouponID:     0,
			//CouponName:   "",
		}
	default:
		return nil, "", errors.New("Invalid order params")

	}
	//折扣码处理
	if preOrder.CouponName != "" {
		msg = couponService.VerifyCoupon(&preOrder)
	}
	return &preOrder, msg, nil
}

// 订单预创建，生成系统订单
func (o *Order) CreateOrder(order *model.Order) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Order{}).Create(&order).Error
	})
}

// 更新数据库订单
func (o *Order) UpdateOrder(order *model.Order) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&order).Error
	})
}

func (o *Order) FirstUserOrder(orderParams *model.Order) (*model.Order, error) {
	var order model.Order
	err := global.DB.Where(&orderParams).First(&order).Error
	return &order, err

}

// 获取用户订单列表
func (o *Order) GetUserOrders(params *model.QueryParams, uID int64) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var orderList []model.Order
	_, dataSql := common_logic.CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	//拼接查询参数
	dataSql = fmt.Sprintf("user_id = %d AND %s", uID, dataSql)
	err := global.DB.Model(&model.Order{}).Count(&data.Total).Where(dataSql).Find(&orderList).Error
	if err != nil {
		return nil, err
	}
	data.Data = orderList
	return &data, nil

}

// 支付成功后，处理普通商品类型的订单
func (o *Order) GoodsTypeGeneralOrderHandler(order *model.Order) error {
	//查询商品信息
	goods, err := shopService.FirstGoods(&model.Goods{ID: order.GoodsID})
	if err != nil {
		return err
	}
	switch goods.DeliverType {
	case constant.DELIVER_TYPE_NONE:
	case constant.DELIVER_TYPE_MANUAL:
	case constant.DELIVER_TYPE_AUTO:
		order.DeliverText = goods.DeliverText
	}
	return o.UpdateOrder(order) //处理发货内容，更新数据库订单状态
}

// 支付成功后，处理订阅类型的订单
func (o *Order) GoodsTypeSubscribeOrderHandler(order *model.Order) error {
	switch order.OrderType {
	case constant.ORDER_TYPE_NEW:
		// 查询商品
		goods, err := shopService.FirstGoods(&model.Goods{ID: order.GoodsID})
		if err != nil {
			return err
		}
		if !goods.IsSale {
			return errors.New(constant.ERROR_GOODS_NOT_SALE)
		}
		// 创建订阅服务
		err = customerService.CreateCustomerService(goods, order)
		if err != nil {
			return err
		}
		//更新数据库订单状态
		return o.UpdateOrder(order)

	case constant.ORDER_TYPE_RENEW:
		// 查找用户服务
		cs, err := customerService.FirstCustomerService(&model.CustomerService{UserID: order.UserID, ID: order.CustomerServiceID})
		if err != nil {
			return err
		}
		// 更新客户服务
		cs.ServiceStatus = true
		//如果没到期，就追加有效期，否则从当天开始设置开始时间
		if cs.ServiceStatus {
			cs.ServiceEndAt = cs.ServiceEndAt.AddDate(0, int(cs.Duration), 0)
		} else {
			cs.ServiceEndAt = time.Now().AddDate(0, int(cs.Duration), 0)
		}
		err = customerService.SaveCustomerService(cs)
		if err != nil {
			return err
		}
		//更新数据库订单状态
		return o.UpdateOrder(order)
	default:
		return errors.New(constant.ERROR_INVALID_ORDER_TYPE)
	}

}

// 支付成功后，处理充值卡类型的订单
func (o *Order) GoodsTypeRechargeOrderHandler(order *model.Order) error {
	err := o.UpdateOrder(order) //更新数据库订单状态
	if err != nil {
		return err
	}
	return userService.RechargeHandle(order) //处理用户余额
}

// 支付成功后，处理订单
func (o *Order) PaymentSuccessfullyOrderHandler(order *model.Order) error {
	var err error
	switch order.GoodsType {
	case constant.GOODS_TYPE_GENERAL: //普通商品
		err = o.GoodsTypeGeneralOrderHandler(order)
	case constant.GOODS_TYPE_SUBSCRIBE: //订阅
		err = o.GoodsTypeSubscribeOrderHandler(order)
	case constant.GOODS_TYPE_RECHARGE: //充值
		err = o.GoodsTypeRechargeOrderHandler(order)
	default:
		err = errors.New(constant.ERROR_INVALID_GOODS_TYPE)
	}
	if err != nil {
		return err
	}
	o.DeleteOneOrderFromCache(order)
	return nil
}

// 删除缓存中的订单
func (o *Order) DeleteOneOrderFromCache(orderParams *model.Order) {
	_, ok := global.LocalCache.Get(constant.CACHE_SUBMIT_ORDER_BY_ORDERID + orderParams.OutTradeNo)
	if ok {
		global.LocalCache.Set(constant.CACHE_SUBMIT_ORDER_BY_ORDERID+orderParams.OutTradeNo,
			&model.Order{},
			5*time.Second)
	}
}

// 更新缓存中的订单
func (o *Order) UpdateOneOrderToCache(orderParams *model.Order) {
	_, ok := global.LocalCache.Get(constant.CACHE_SUBMIT_ORDER_BY_ORDERID + orderParams.OutTradeNo)
	if ok {
		global.LocalCache.
			Set(constant.CACHE_SUBMIT_ORDER_BY_ORDERID+orderParams.OutTradeNo,
				orderParams,
				constant.CACHE_SUBMIT_ORDER_TIMEOUT*time.Minute)
	}
}

// 处理超时订单
func (o *Order) OrderTimeout(k string, v any) {
	if strings.Index(k, constant.CACHE_SUBMIT_ORDER_BY_ORDERID) != -1 {
		preOrder := v.(*model.Order)
		switch preOrder.OrderType {
		case constant.ORDER_TYPE_NEW:
		case constant.ORDER_TYPE_RENEW:
		default:
			return
		}
		// 关闭订单
		preOrder.OrderType = constant.ORDER_TYPE_RESTORE
		global.Queue.Publish(constant.SUBMIT_ORDER, preOrder)
	}
}
