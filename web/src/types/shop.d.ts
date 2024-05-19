declare interface Goods {
  created_at: string;
  updated_at: string;
  id: number;
  good_order: number;    //排序
  cover_image:string     //封面图片
  subject: string;       // 标题
  des: string;     //商品描述
  price: string;  // 商品价格
  is_show: boolean; //是否显示
  is_sale: boolean;//是否售卖，是否上架
  is_renew: boolean;//是否可续费
  quota: number;//限购数量
  stock: number;//库存
  goods_type: string;    //商品类型，general=普通商品 subscribe=订阅 recharge=充值
  deliver_type: string;  //发货类型，none=不发货，auto=自动发货，manual=手动发货
  deliver_text: string;  //发货内容;type:text

  //订阅参数
  price_3_month:string
  price_6_month:string
  price_12_month:string
  price_unlimited_duration:string
  enable_traffic_reset:boolean
  total_bandwidth: number; //总流量
  node_connector: number; //可连接客户端数量
  node_speed_limit: number;// 限速Mbps（Mbps）

  //充值参数
  recharge_amount: string;      //充值金额
  nodes: NodeInfo[];
}

declare interface PayInfo {
  created_at: string;
  updated_at: string;
  id: number;
  name: string;
  pay_type: string;
  pay_logo_url: string;
  status: boolean;
  alipay: Alipay;
  epay: Epay;
}

declare interface Alipay {
  alipay_app_id: string;
  alipay_notify_url: string;
  alipay_app_private_key: string;
  alipay_ali_public_key: string;
  alipay_encrypt_key: string;

}

declare interface Epay {
  epay_pid: number;
  epay_key: string;
  epay_api_url: string;
  epay_return_url: string;
  epay_notify_url: string;
  epay_type: string;

}

// 易支付支付结果响应
declare interface EpayResultResponse {
  pid: number;
  trade_no: string;
  out_trade_no: string;
  type: string;
  name: string;
  money: string;
  tradeStatus: string;
  param: string;
  sign: string;
  sign_type: string;
}

declare interface PreCreatePayToFrontend {
  alipay_info: AlipayPreCreatePayToFrontend;
  epay_info: EpayPreCreatePayToFrontend;
}

declare interface AlipayPreCreatePayToFrontend {
  qr_code: string;
}

// 易支付支付预创建返回给前端
declare interface EpayPreCreatePayToFrontend {
  epay_api_url: string;
  epay_pre_create_pay: EpayPreCreatePay;

}

declare interface EpayPreCreatePay {
  pid: number;
  type: string;
  out_trade_no: string;
  notify_url: string;
  return_url: string;
  name: string;
  money: string;
  clientip: string;
  device: string;
  param: string;
  sign: string;
  sign_type: string;
}