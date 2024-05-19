export default {
  adminOrder: {
    query:"Query",
    modify_order:"Modify order",
   submitOrder:"Submit order",
   orderDetails:"Order Details",
     Order: {
      index:"index",
      created_at: "created at",
      updated_at: "updated at",
      id: "id",
      order_type:"order type", //订单类型:New=新购入;Renew=续费"
      trade_status: "trade status", //交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"
      out_trade_no: "out trade no", //商户订单号
      original_amount: "original amount", //原始金额
      total_amount: "total amount", //订单金额
      buyer_pay_amount: "buyer pay amount", //付款金额
      coupon_amount: "coupon amount", //折扣码折扣金额
      balance_amount: "balance amount", //余额折扣金额
      // 关联用户
      user_id: "user id",
      user_name: "user name",
      // user: "",
      // 商品参数
      goods_id: "goods id",
      des: "des", //描述
      goods_type: "goods type",  //类型
      deliver_type: "deliver type", //发货类型
      deliver_text: "deliver text", //发货内容
      subject: "subject",
      price: "price",
      duration: "duration", //购买时长(单位：月)
      // 服务参数
      customer_service_id:"customer service id",
      //支付参数
      pay_id: "pay id",   //支付方式id
      pay_type: "pay type", //支付类型，alipay,epay
      coupon_id: "coupon id",
      coupon_name: "coupon name",
      // pay_info: "", //支付信息，epay，alipay等
      trade_no: "trade no",
      buyer_logon_id: "buyer logon id",
    }
  },
};