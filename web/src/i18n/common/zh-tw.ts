export default {
  common: {
    tip:"提示",
    message_confirm_delete:"此操作將永久删除數據，是否繼續?",
    button_cancel:"取消",
    button_confirm:"確認",

    operate:"操作",
    delete:"删除",
    modify:"修改",
    reset:"重置",
    query:"査詢",
    advanced_query:"高級查詢",
    add:"新增",
    sort:"排序",
    reply:"回復",
    close:"關閉",
    back:"上一步",
    send:"發送",
    details:"詳情",
    purchase:"購買",

    yes:"是",
    no:"否",
    enable:"啟用",
    disable:"禁用",
    display:"顯示",
    hide:"隱藏",

    startDate:"開始日期",
    endDate:"結束日期",
    to:"至",
    up:"起",
    lastWeek:"上周",
    lastMonth:"上月",
    lastThreeMonths:"最近3個月",
    month:"月",
    no_time_limit:"不限時",

    noData:"暫無數據",

    success:"成功",
    warning:"警告",
    error:"錯誤",
  },
  constant:{
    // 商品
    GOODS_TYPE_GENERAL   : "普通商品",
    GOODS_TYPE_SUBSCRIBE : "訂閱",
    GOODS_TYPE_RECHARGE  : "充值",
    // 订单类型
    ORDER_TYPE_NEW:"新購",
    ORDER_TYPE_RENEW:"續費",
    ORDER_TYPE_DESTROYED:"已銷毀",
    // 订单状态
    ORDER_STATUS_CREATED        : "已創建",
    ORDER_STATUS_WAIT_BUYER_PAY : "等待付款",
    ORDER_STATUS_TRADE_CLOSED   : "交易關閉",
    ORDER_STATUS_TRADE_SUCCESS  : "交易成功",
    ORDER_STATUS_TRADE_FINISHED : "交易结束",
    ORDER_STATUS_COMPLETED      : "交易完成 ",
    ORDER_STATUS_UNKNOWN_STATE  : "未知狀態 ",
    // 发货
    DELIVER_TYPE_NONE   : "無",
    DELIVER_TYPE_AUTO   : "自動",
    DELIVER_TYPE_MANUAL : "手動",
    //工单状态
    TICKET_PROCESSING : "工單進行中",//工单进行中
    TICKET_CLOSED     :"工單已結束",
    //邮件类型
    EMAIL_TYPE_USER_REGISTER : "EMAIL_TYPE_USER_REGISTER",
    EMAIL_TYPE_USER_RESETPWD : "EMAIL_TYPE_USER_RESETPWD",
    //节点类型
    NODE_TYPE_NORMAL    : "正常节点",
    NODE_TYPE_SHARED    : "共亯節點",
    NODE_TYPE_TRANSFER    : "中轉節點",
    //节点协议类型
    NODE_PROTOCOL_VMESS       : "vmess",
    NODE_PROTOCOL_VLESS       : "vless",
    NODE_PROTOCOL_TROJAN      : "trojan",
    NODE_PROTOCOL_HYSTERIA    : "hysteria2",
    NODE_PROTOCOL_SHADOWSOCKS : "shadowsocks",
    //文章类型
    ARTICLE_TYPE_HOME      : "首頁",
    ARTICLE_TYPE_DIALOG    : "彈窗",
    ARTICLE_TYPE_NOTICE    : "公告",
  },
};