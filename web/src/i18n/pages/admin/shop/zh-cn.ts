export default {
  adminShop:{
    addGoods:"新增商品",
    modifyGoods:"修改商品",
    selectGoods:"选择商品",
    sort:"排序",
    purchase:"购买",
    resultText1:"请返回首页查看服务状态",
    resultText2:"请尽快完成支付",
    resultText3:"支付遇到问题",
    resultText4:"前往支付宝",
    resultText5:"前往支付",
    resultText6:"请求超时，请稍后检查我的订单，有无订单记录",
    Goods:{
      index:"序号",
      created_at: "创建时间",
      updated_at: "更新时间",
      id: "id",
      good_order: "排序",    //排序
      subject: "标题",       // 标题
      cover_image:"封面图片",//封面图片
      des: "商品描述",     //商品描述
      price: "价格",  // 商品价格
      price_3_month:"3个月的价格",
      price_6_month:"6个月的价格",
      price_12_month:"12个月的价格",
      price_unlimited_duration:"不限时的价格",
      is_show: "是否显示", //是否显示
      is_sale: "是否售卖",//是否售卖
      is_renew: "是否可续费",//是否可续费
      quota: "限购数量",//限购数量
      stock: "库存",//库存
      goods_type: "商品类型",    //类型，general=普通商品 subscribe=订阅 recharge=充值
      deliver_type: "发货类型",  //发货类型，none=不发货，auto=自动发货，manual=手动发货
      deliver_text: "发货内容",  //发货内容;type:text
      //订阅参数
      enable_traffic_reset:"自动重置流量",
      total_bandwidth: "总流量", //总流量
      node_connector: "节点连接数", //可连接客户端数量
      node_speed_limit: "限速",// 限速Mbps（Mbps）
      //充值参数
      recharge_amount: "充值金额",      //充值金额
      nodes: "关联节点",
    },
    PayInfo:{
      index:"序号",
      created_at: "创建时间",
      updated_at: "更新时间",
      id: "id",
      name: "名称",
      pay_type: "支付类型",
      pay_logo_url: "logo",
      status: "状态",
      alipay: "alipay",
      epay: "epay",
    },
    Alipay:{
      alipay_app_id: "支付宝appID",
      alipay_notify_url: "支付宝异步回调地址",
      alipay_app_private_key: "支付宝应用私钥",
      alipay_ali_public_key: "支付宝公钥",
      alipay_encrypt_key: "支付宝接口加密密钥",
    },
    Epay:{
      epay_pid: "易支付商户ID",
      epay_key: "易支付商户密钥",
      epay_main_url: "易支付网址",
      epay_api_url: "易支付api地址",
      epay_return_url: "易支付跳转通知地址",
      epay_notify_url: "易支付异步通知地址",
      epay_type: "支付类型",
    },
    //
    coupon:"折扣码",
    couponRateTip:"价格 = 原价 * (1- 折扣率)",
    addCoupon:"新增折扣",
    modifyCoupon:"修改折扣",
    Coupon:{
      index:"序号",
      created_at: "创建时间",
      updated_at: "更新时间",
      id:"id",
      name: "折扣码名称",
      discount_rate:"折扣率",
      limit:"次数",
      expired_at: "过期时间",
      min_amount:"最低使用金额",
      goods:"关联商品",
    }
  }
}
