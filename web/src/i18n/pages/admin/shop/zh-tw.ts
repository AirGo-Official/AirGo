export default {
  adminShop:{
    addGoods:"新增商品",
    modifyGoods:"修改商品",
    selectGoods:"選擇商品",
    sort:"排序",
    purchase:"購買",
    resultText1:"請返回首頁查看服務狀態",
    resultText2:"請儘快完成支付",
    resultText3:"支付遇到問題",
    resultText4:"前往支付寶",
    resultText5:"前往支付",
    resultText6:"請求超時，請稍後檢查我的訂單，有無訂單記錄",
    Goods:{
      index:"序號",
      created_at: "創建時間",
      updated_at: "更新時間",
      id: "id",
      good_order: "排序",    //排序
      subject: "標題",       // 标题
      cover_image:"封面图片",//封面图片
      des: "商品描述",     //商品描述
      price: "價格",  // 商品价格
      price_3_month:"3個月的價格",
      price_6_month:"6個月的價格",
      price_12_month:"12個月的價格",
      price_unlimited_duration:"不限時的價格",
      is_show: "是否顯示", //是否显示
      is_sale: "是否售賣",//是否售卖
      is_renew: "是否可續費",//是否可续费
      quota: "限购数量",//限购数量
      stock: "庫存",//库存
      goods_type: "商品類型",    //类型，general=普通商品 subscribe=订阅 recharge=充值
      deliver_type: "發貨類型",  //发货类型，none=不发货，auto=自动发货，manual=手动发货
      deliver_text: "發貨內容",  //发货内容;type:text
      //订阅参数
      enable_traffic_reset:"自動重置流量",
      total_bandwidth: "總流量", //总流量
      node_connector: "節點連接數", //可连接客户端数量
      node_speed_limit: "限速",// 限速Mbps（Mbps）
      //充值参数
      recharge_amount: "充值金額",      //充值金额
      nodes: "關聯節點",
    },
    PayInfo:{
      index:"序號",
      created_at: "創建時間",
      updated_at: "更新時間",
      id: "id",
      name: "名稱",
      pay_type: "支付類型",
      pay_logo_url: "logo",
      status: "狀態",
      alipay: "alipay",
      epay: "epay",
    },
    Alipay:{
      alipay_app_id: "支付寶appID",
      alipay_notify_url: "支付寶非同步回檔地址",
      alipay_app_private_key: "支付寶應用私密金鑰",
      alipay_ali_public_key: "支付寶公開金鑰",
      alipay_encrypt_key: "支付寶介面加密金鑰",
    },
    Epay:{
      epay_pid: "易支付商戶ID",
      epay_key: "易支付商戶金鑰",
      epay_main_url: "易支付網址",
      epay_api_url: "易支付api地址",
      epay_return_url: "易支付跳轉通知地址",
      epay_notify_url: "易支付非同步通知地址",
      epay_type: "支付類型",
    },
    //
    coupon:"折扣碼",
    couponRateTip:"價格=原價*（1-折扣率）",
    addCoupon:"新增折扣",
    modifyCoupon:"修改折扣",
    Coupon:{
      index:"序號",
      created_at: "創建時間",
      updated_at: "更新時間",
      id:"id",
      name: "折扣碼名稱",
      discount_rate:"折扣率",
      limit:"次數",
      expired_at: "過期時間",
      min_amount:"最低使用金額",
      goods:"關聯商品",
    }
  }
}
