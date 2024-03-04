declare interface Order {
    created_at: string
    updated_at: string
    id: number
    order_type:string //订单类型:New=新购入;Renew=续费"
    trade_status: string //交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"
    out_trade_no: string //商户订单号
    original_amount: string //原始金额
    total_amount: string //订单金额
    buyer_pay_amount: string //付款金额
    coupon_amount: string //折扣码折扣金额
    balance_amount: string //余额折扣金额
    // 关联用户
    user_id: number
    user_name: string
    user: any
    // 商品参数
    goods_id: number
    des: string //描述
    goods_type: string  //类型
    deliver_type: string //发货类型
    deliver_text: string //发货内容
    subject: string
    price: string
    duration: number //购买时长(单位：月)

    // 服务参数
    customer_service_id:number

    //支付参数
    pay_id: number   //支付方式id
    pay_type: string //支付类型，alipay,epay
    coupon_id: number
    coupon_name: string
    pay_info: PreCreatePayToFrontend; //支付信息，epay，alipay等
    trade_no: string
    buyer_logon_id: string

}

declare interface OrdersWithTotal {
    total_amount: number;
    total: number;
}
declare interface OrderSummary {
    date:string
    order_total:number
    income_total:number
    general_total:number
    recharge_total:number
    subscribe_total:number
}