export default {
  adminOrder: {
    query:"査詢",
    modify_order:"修改訂單",
    submitOrder:"提交訂單",
    orderDetails:"订单详情",
    Order: {
      index:"序號",
      created_at: "創建時間",
      updated_at: "更新時間",
      id: "id",
      order_type:"訂單類型", //订单类型:New=新购入;Renew=续费"
      trade_status: "交易狀態", //交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"
      out_trade_no: "商戶訂單號", //商户订单号
      original_amount: "原始金額", //原始金额
      total_amount: "訂單金額", //订单金额
      buyer_pay_amount: "付款金額", //付款金额
      coupon_amount: "折扣碼折扣金額", //折扣码折扣金额
      balance_amount: "餘額折扣金額", //余额折扣金额
      // 关联用户
      user_id: "用戶id",
      user_name: "用戶名",
      // user: "",
      // 商品参数
      goods_id: "商品id",
      des: "des", //描述
      goods_type: "商品類型",  //类型
      deliver_type: "發貨類型", //发货类型
      deliver_text: "發貨內容", //发货内容
      subject: "商品標題",
      price: "商品價格",
      duration: "訂購時長", //购买时长(单位：月)
      // 服务参数
      customer_service_id:"服務id",
      //支付参数
      pay_id: "支付id",   //支付方式id
      pay_type: "支付類型", //支付方式，alipay,epay
      coupon_id: "折扣碼id",
      coupon_name: "折扣碼內容",
      // pay_info: "", //支付信息，epay，alipay等
      trade_no: "協力廠商訂單號",
      buyer_logon_id: "協力廠商用戶id",
    }
  },
};