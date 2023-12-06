package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// 获取全部订单，分页获取
func GetAllOrder(ctx *gin.Context) {
	var params model.FieldParamsReq
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllOrder error:"+err.Error(), nil, ctx)
		return
	}
	res, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllOrder error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetAllOrder success", model.CommonDataResp{
		Total: total,
		Data:  res,
	}, ctx)

}

// 获取订单统计
func GetMonthOrderStatistics(ctx *gin.Context) {
	var params model.FieldParamsReq
	err := ctx.ShouldBind(&params)
	res, err := service.GetMonthOrderStatistics(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetMonthOrderStatistics error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetMonthOrderStatistics success", res, ctx)
}

// 获取用户订单by user id
func GetOrderByUserID(ctx *gin.Context) {
	var params model.FieldParamsReq
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllOrder error:"+err.Error(), nil, ctx)
		return
	}

	uIDInt, ok := GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("GetOrderByUserID error:user id error", nil, ctx)
		return
	}
	res, err := service.GetUserOrders(&params, uIDInt)
	if err != nil {
		response.Fail("GetOrderByUserID error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetOrderByUserID success", res, ctx)
}

// 完成未支付订单
func CompletedOrder(ctx *gin.Context) {
	var order model.Orders
	err := ctx.ShouldBind(&order)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("CompletedOrder error:"+err.Error(), nil, ctx)
		return
	}
	order.TradeStatus = model.OrderCompleted //更新数据库订单状态,自定义结束状态Completed
	err = service.UpdateOrder(&order)        //更新数据库状态
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("CompletedOrder error:"+err.Error(), nil, ctx)
		return
	}
	err = service.UpdateUserSubscribe(&order) //更新用户订阅信息
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("CompletedOrder error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("CompletedOrder success", nil, ctx)

}

// 获取订单详情（计算价格等）
func GetOrderInfo(ctx *gin.Context) {
	order, msg := PreHandleOrder(ctx)
	if order == nil {
		response.Fail("GetOrderInfo error:order is null", nil, ctx)
		return
	}
	if msg == "" {
		msg = "GetOrderInfo success"
	}
	response.OK(msg, order, ctx)
}

// 订单预创建，生成系统订单
func PreCreateOrder(ctx *gin.Context) {
	order, _ := PreHandleOrder(ctx)
	if order == nil {
		response.Fail("PreCreateOrder error:order is null", nil, ctx)
		return
	}
	//创建系统订单
	order.TradeStatus = model.OrderCreated
	err := service.CommonSqlCreate[model.Orders](*order)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreCreateOrder error::"+err.Error(), nil, ctx)
		return
	}
	response.OK("PreCreateOrder success", order, ctx)
}

// 订单预处理，计算价格
func PreHandleOrder(ctx *gin.Context) (*model.Orders, string) {

	uIDInt, _ := GetUserIDFromGinContext(ctx)
	uName, _ := GetUserNameFromGinContext(ctx)

	var msg string
	user, _ := service.FindUserByID(uIDInt)

	var receiveOrder model.Orders
	err := ctx.ShouldBind(&receiveOrder) //前端传过来 goods_id,coupon_name
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreHandleOrder error:"+err.Error(), nil, ctx)
		return nil, ""
	}
	//通过商品id查找商品
	goods, err := service.FindGoodsByGoodsID(receiveOrder.GoodsID)
	if err != nil {
		global.Logrus.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ""
		}
	}
	//构造系统订单参数
	uIDStr := other_plugin.Sup(uIDInt, 6) //对长度不足n的后面补0
	receiveOrder.GoodsID = goods.ID       //商品ID

	receiveOrder.OutTradeNo = time.Now().Format("20060102150405") + uIDStr //系统订单号：时间戳+6位user id
	receiveOrder.Subject = goods.Subject                                   //商品的标题
	receiveOrder.Price = goods.TotalAmount                                 //商品的价格
	receiveOrder.TotalAmount = goods.TotalAmount                           //订单的价格
	receiveOrder.UserID = uIDInt                                           //用户ID
	receiveOrder.UserName = uName                                          //用户名
	//receiveOrder.PayType = receiveOrder.PayType //添加付款方式
	//折扣码处理
	total, _ := strconv.ParseFloat(goods.TotalAmount, 64)
	if receiveOrder.CouponName != "" {
		coupon, err := service.VerifyCoupon(&receiveOrder)
		if err != nil {
			global.Logrus.Error(err.Error())
			msg = err.Error()
		}
		if coupon.DiscountRate != 0 {
			receiveOrder.CouponAmount = fmt.Sprintf("%.2f", total*coupon.DiscountRate)
			receiveOrder.CouponID = coupon.ID
			total = total - total*coupon.DiscountRate //total-折扣码
		}
	}
	//旧套餐抵扣处理
	if global.Server.Subscribe.EnabledDeduction {
		//计算剩余率
		if user.SubscribeInfo.SubStatus {
			rate, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64((user.SubscribeInfo.T-user.SubscribeInfo.U-user.SubscribeInfo.D))/float64(user.SubscribeInfo.T)), 64)
			//if math.IsNaN(rate) {
			if err != nil {
				rate = 0 //
			}
			//套餐流量剩余率大于设定的阈值才进行处理
			if rate >= global.Server.Subscribe.DeductionThreshold {
				//查找旧套餐价格
				order, _, _ := service.CommonSqlFind[model.Orders, string, model.Orders](fmt.Sprintf("user_id = %d ORDER BY id desc LIMIT 1", uIDInt))
				if order.ReceiptAmount != "" { //使用 实收金额 进行判断
					receiptAmount, _ := strconv.ParseFloat(order.ReceiptAmount, 64)
					deductionAmount := receiptAmount * rate
					if deductionAmount < total {
						receiveOrder.DeductionAmount = fmt.Sprintf("%.2f", receiptAmount*rate)
						total = total - deductionAmount
					} else {
						receiveOrder.DeductionAmount = fmt.Sprintf("%.2f", total)
						total = 0
					}
				}
			}
		}
	}
	//余额抵扣，计算最终价格，TotalAmount=总价-折扣码的折扣-旧套餐的抵扣
	if user.Remain > 0 {
		if user.Remain < total {
			receiveOrder.RemainAmount = fmt.Sprintf("%.2f", user.Remain)
			total = total - user.Remain
		} else {
			receiveOrder.RemainAmount = fmt.Sprintf("%.2f", total)
			total = 0
		}
	}
	receiveOrder.TotalAmount = fmt.Sprintf("%.2f", total)
	return &receiveOrder, msg
}
