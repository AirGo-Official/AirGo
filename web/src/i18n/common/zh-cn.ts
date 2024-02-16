export default {
  common: {
    tip:"提示",
    message_confirm_delete:"此操作将永久删除数据，是否继续?",

    button_cancel:"取消",
    button_confirm:"确认",

    operate:"操作",
    delete:"删除",
    modify:"修改",
    reset:"重置",


    yes:"是",
    no:"否",
    enable:"启用",
    disable:"禁用",

    query:"查询",
    advanced_query:"高级查询",
    add:"新增",
    sort:"排序",

    startDate:"开始日期",
    endDate:"结束日期",
    to:"至",
    lastWeek:"上周",
    lastMonth:"上月",
    lastThreeMonths:"最近3个月",
  },
  constant:{
    // 商品
    GOODS_TYPE_GENERAL   : "普通商品",
    GOODS_TYPE_SUBSCRIBE : "订阅",
    GOODS_TYPE_RECHARGE  : "充值",
    // 订单类型
    ORDER_TYPE_NEW:"新购",
    ORDER_TYPE_RENEW:"续费",
    // 订单状态
    ORDER_STATUS_CREATED        : "已创建",
    ORDER_STATUS_WAIT_BUYER_PAY : "等待付款",
    ORDER_STATUS_TRADE_CLOSED   : "交易关闭",
    ORDER_STATUS_TRADE_SUCCESS  : "交易成功",
    ORDER_STATUS_TRADE_FINISHED : "交易结束",
    ORDER_STATUS_COMPLETED      : "交易完成 ",
    // 发货
    DELIVER_TYPE_NONE   : "无",
    DELIVER_TYPE_AUTO   : "自动",
    DELIVER_TYPE_MANUAL : "手动",
    //工单状态
    TICKET_PROCESSING : "工单进行中",//工单进行中
    TICKET_CLOSED     :"工单已结束",
    //邮件类型
    EMAIL_TYPE_USER_REGISTER : "EMAIL_TYPE_USER_REGISTER",
    EMAIL_TYPE_USER_RESETPWD : "EMAIL_TYPE_USER_RESETPWD",
    //节点类型
    NODE_TYPE_TRANSFER    : "transfer",
    NODE_TYPE_VMESS       : "vmess",
    NODE_TYPE_VLESS       : "vless",
    NODE_TYPE_TROJAN      : "trojan",
    NODE_TYPE_HYSTERIA    : "hysteria",
    NODE_TYPE_SHADOWSOCKS : "shadowsocks",
  },
};