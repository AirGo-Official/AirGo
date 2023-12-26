declare interface Order {
    created_at: string;
    updated_at: string;
    id: number;
    user_id: number;
    user_name: string;
    user: any;

    out_trade_no: string;
    goods_id: number;
    goods_type: string;  //类型
    deliver_type: string;//发货类型
    deliver_text: string;//发货内容
    subject: string;
    price: string;
    pay_id: number;   //支付方式id
    pay_type: string; //支付方式，alipay,epay
    coupon_id: number;
    coupon_name: string;
    coupon_amount: string;
    deduction_amount: string;
    remain_amount: string;

    pay_info: PreCreatePayToFrontend; //支付信息，epay，alipay等
    trade_no: string;
    buyer_logon_id: string;
    trade_status: string;
    total_amount: string;
    receipt_amount: string;
    buyer_pay_amount: string;
}

declare interface OrdersWithTotal {
    total_amount: number;
    total: number;
}