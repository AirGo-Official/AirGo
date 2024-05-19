export default {
  adminShop:{
    addGoods:"Add goods",
    modifyGoods:"Modify goods",
    selectGoods:"Select goods",
    sort:"Sort",
    purchase:"Purchase",
    resultText1:"Please return to the homepage to check the service status",
    resultText2:"Please complete payment as soon as possible",
    resultText3:"Please complete the payment as soon as possible",
    resultText4:"Go to Alipay",
    resultText5:"Go to Pay",
    resultText6:"Request timeout, please check my order later for any order records",
    Goods:{
      index:"index",
      created_at: "created at",
      updated_at: "updated at",
      id: "id",
      good_order: "good order",    //排序
      cover_image:"cover image",//封面图片
      subject: "subject",       // 标题
      des: "des",     //商品描述
      price: "price",  // 商品价格
      price_3_month:"price for 3 months",
      price_6_month:"price for 6 months",
      price_12_month:"price for 12 months",
      price_unlimited_duration:"price for unlimited",
      is_show: "is show", //是否显示
      is_sale: "is sale",//是否售卖，是否上架
      is_renew: "is renew",//是否可续费
      quota: "quota",//限购数量
      stock: "stock",//库存
      goods_type: "goods type",    //类型，general=普通商品 subscribe=订阅 recharge=充值
      deliver_type: "deliver type",  //发货类型，none=不发货，auto=自动发货，manual=手动发货
      deliver_text: "deliver text",  //发货内容;type:text
      //订阅参数
      enable_traffic_reset:"enable traffic reset",
      total_bandwidth: "total bandwidth", //总流量
      node_connector: "node connector", //可连接客户端数量
      node_speed_limit: "node speed limit",// 限速Mbps（Mbps）
      //充值参数
      recharge_amount: "recharge amount",  //充值金额
      nodes: "nodes",
    },
    PayInfo:{
      index:"index",
      created_at: "created at",
      updated_at: "updated at",
      id: "id",
      name: "name",
      pay_type: "pay type",
      pay_logo_url: "pay logo url",
      status: "status",
      alipay: "alipay",
      epay: "epay",
    },
    Alipay:{
      alipay_app_id: "alipay app id",
      alipay_notify_url: "alipay notify url",
      alipay_app_private_key: "alipay app private key",
      alipay_ali_public_key: "alipay ali public key",
      alipay_encrypt_key: "alipay encrypt key",
    },
    Epay:{
      epay_pid: "epay pid",
      epay_key: "epay key",
      epay_api_url: "epay api url",
      epay_main_url: "epay main url",
      epay_return_url: "epay return url",
      epay_notify_url: "epay notify url",
      epay_type: "epay type",
    },
    //
    coupon:"Coupon",
    couponRateTip:"price = price * (1- rate)",
    addCoupon:"Add coupon",
    modifyCoupon:"Modify coupon",
    Coupon:{
      index:"index",
      created_at: "created at",
      updated_at: "updated at",
      id:"id",
      name: "coupon name",
      discount_rate:"discount rate",
      limit:"limit",
      expired_at: "expired at",
      min_amount:"min amount",
      goods:"goods",
    }
  }
}
