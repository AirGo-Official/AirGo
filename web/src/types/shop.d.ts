declare interface Goods {
    id: number;
    created_at: string;
    updated_at: string;
    //
    subject: string;       // 标题
    total_amount: string;  // 总金额
    good_order: number;    //排序
    status: boolean; //是否启用
    des: string;     //商品描述
    goods_type: string;    //类型，general=普通商品 subscribe=订阅 recharge=充值
    deliver_type: string;  //发货类型，none=不发货，auto=自动发货，manual=手动发货
    deliver_text: string;  //发货内容;type:text
    //订阅参数
    total_bandwidth: number; //总流量
    expiration_date: number; //有效期
    traffic_reset_method: string; //流量重置方式,Stack,NotStack
    reset_day: number;      //流量重置日
    node_connector: number; //可连接客户端数量
    //充值参数
    recharge_amount: string;      //充值金额
    //
    checked_nodes: number[]; //套餐编辑时选中的节点
    nodes: Node[];
}

declare interface PayInfo {
    created_at: string
    updated_at: string
    id: number
    name: string
    pay_type: string
    pay_logo_url: string
    status: boolean
    alipay: Alipay
    epay: Epay
}

declare interface Alipay {
    alipay_app_id: string
    alipay_return_url: string
    alipay_app_private_key: string
    alipay_ali_public_key: string
    alipay_encrypt_key: string

}

declare interface Epay {
    epay_pid: number
    epay_key: string
    epay_api_url: string
    epay_return_url: string
    epay_notify_url: string
    epay_type: string

}

// 易支付支付结果响应
declare interface EpayResultResponse {
    pid: number
    trade_no: string
    out_trade_no: string
    type: string
    name: string
    money: string
    tradeStatus: string
    param: string
    sign: string
    sign_type: string
}

declare interface PreCreatePayToFrontend {
    alipay_info: AlipayPreCreatePayToFrontend
    epay_info: EpayPreCreatePayToFrontend
}

declare interface AlipayPreCreatePayToFrontend {
    qr_code: string
}

// 易支付支付预创建返回给前端
declare interface EpayPreCreatePayToFrontend {
    epay_api_url: string
    epay_pre_create_pay: EpayPreCreatePay

}

declare interface EpayPreCreatePay {
    pid: number
    type: string
    out_trade_no: string
    notify_url: string
    return_url: string
    name: string
    money: string
    clientip: string
    device: string
    param: string
    sign: string
    sign_type: string
}