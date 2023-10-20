declare interface Goods {
    id: number;
    created_at: string;
    updated_at: string;

    good_order: number;    //排序
    status: boolean; //是否启用
    des: string;     //商品描述
    subject: string;       // 标题
    total_amount: string;  // 总金额
    total_bandwidth: number; //总流量
    expiration_date: number; //有效期
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