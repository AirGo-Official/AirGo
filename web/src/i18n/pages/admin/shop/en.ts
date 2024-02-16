export default {
  adminShop:{
    addGoods:"Add goods",
    modifyGoods:"Modify goods",
    sort:"Sort",
    Goods:{
      index:"index",
      created_at: "created_at",
      updated_at: "updated_at",
      id: "id",
      good_order: "good_order",    //排序
      subject: "subject",       // 标题
      des: "des",     //商品描述
      price: "price",  // 商品价格
      is_show: "is_show", //是否显示
      is_sale: "is_sale",//是否售卖，是否上架
      is_renew: "is_renew",//是否可续费
      quota: "quota",//限购数量
      stock: "stock",//库存
      goods_type: "goods_type",    //类型，general=普通商品 subscribe=订阅 recharge=充值
      deliver_type: "deliver_type",  //发货类型，none=不发货，auto=自动发货，manual=手动发货
      deliver_text: "deliver_text",  //发货内容;type:text
      //订阅参数
      total_bandwidth: "total_bandwidth", //总流量
      node_connector: "node_connector", //可连接客户端数量
      node_speed_limit: "node_speed_limit",// 限速Mbps（Mbps）
      //充值参数
      recharge_amount: "recharge_amount",  //充值金额
      nodes: "nodes",
    },
    PayInfo:{
      index:"index",
      created_at: "created_at",
      updated_at: "updated_at",
      id: "id",
      name: "name",
      pay_type: "pay_type",
      pay_logo_url: "pay_logo_url",
      status: "status",
      alipay: "alipay",
      epay: "epay",
    },
    Alipay:{
      alipay_app_id: "alipay_app_id",
      alipay_return_url: "alipay_return_url",
      alipay_app_private_key: "alipay_app_private_key",
      alipay_ali_public_key: "alipay_ali_public_key",
      alipay_encrypt_key: "alipay_encrypt_key",
    },
    Epay:{
      epay_pid: "epay_pid",
      epay_key: "epay_key",
      epay_api_url: "epay_api_url",
      epay_return_url: "epay_return_url",
      epay_notify_url: "epay_notify_url",
      epay_type: "epay_type",
    },
  }
}