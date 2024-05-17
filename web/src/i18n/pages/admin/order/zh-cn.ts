export default {
  adminOrder: {
    query:"查询",
    modify_order:"修改订单",
    submitOrder:"提交订单",
    orderDetails:"订单详情",
    Order: {
      index:"序号",
      created_at: "创建时间",
      updated_at: "更新时间",
      id: "id",
      order_type:"订单类型", //订单类型:New=新购入;Renew=续费"
      trade_status: "交易状态", //交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"
      out_trade_no: "商户订单号", //商户订单号
      original_amount: "原始金额", //原始金额
      total_amount: "订单金额", //订单金额
      buyer_pay_amount: "付款金额", //付款金额
      coupon_amount: "折扣码折扣金额", //折扣码折扣金额
      balance_amount: "余额折扣金额", //余额折扣金额
      // 关联用户
      user_id: "用户id",
      user_name: "用户名",
      // user: "",
      // 商品参数
      goods_id: "商品id",
      des: "des", //描述
      goods_type: "商品类型",  //类型
      deliver_type: "发货类型", //发货类型
      deliver_text: "发货内容", //发货内容
      subject: "商品标题",
      price: "商品价格",
      duration: "订购时长", //购买时长(单位：月)
      // 服务参数
      customer_service_id:"服务id",
      //支付参数
      pay_id: "支付id",   //支付方式id
      pay_type: "支付方式", //支付方式，alipay,epay
      coupon_id: "折扣码id",
      coupon_name: "折扣码内容",
      // pay_info: "", //支付信息，epay，alipay等
      trade_no: "第三方订单号",
      buyer_logon_id: "第三方用户id",
    }
  },
};