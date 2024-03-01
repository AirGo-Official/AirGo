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
    query:"查询",
    advanced_query:"高级查询",
    add:"新增",
    sort:"排序",
    reply:"回复",
    close:"关闭",
    back:"上一步",
    send:"发送",
    details:"详情",
    purchase:"购买",

    yes:"是",
    no:"否",
    enable:"启用",
    disable:"禁用",
    display:"显示",
    hide:"隐藏",

    startDate:"开始日期",
    endDate:"结束日期",
    to:"至",
    lastWeek:"上周",
    lastMonth:"上月",
    lastThreeMonths:"最近3个月",
    month:"月",

    noData:"暂无数据",

    success:"成功",
    warning:"警告",
    error:"错误",
  },
  constant:{
    // 商品
    GOODS_TYPE_GENERAL   : "普通商品",
    GOODS_TYPE_SUBSCRIBE : "订阅",
    GOODS_TYPE_RECHARGE  : "充值",
    // 订单类型
    ORDER_TYPE_NEW:"新购",
    ORDER_TYPE_RENEW:"续费",
    ORDER_TYPE_DESTROYED:"已销毁",
    // 订单状态
    ORDER_STATUS_CREATED        : "已创建",
    ORDER_STATUS_WAIT_BUYER_PAY : "等待付款",
    ORDER_STATUS_TRADE_CLOSED   : "交易关闭",
    ORDER_STATUS_TRADE_SUCCESS  : "交易成功",
    ORDER_STATUS_TRADE_FINISHED : "交易结束",
    ORDER_STATUS_COMPLETED      : "交易完成 ",
    ORDER_STATUS_UNKNOWN_STATE  : "未知状态 ",
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
    NODE_TYPE_NORMAL    : "正常节点",
    NODE_TYPE_SHARED    : "共享节点",
    NODE_TYPE_TRANSFER    : "中转节点",
    //节点协议类型
    NODE_PROTOCOL_VMESS       : "vmess",
    NODE_PROTOCOL_VLESS       : "vless",
    NODE_PROTOCOL_TROJAN      : "trojan",
    NODE_PROTOCOL_HYSTERIA    : "hysteria",
    NODE_PROTOCOL_SHADOWSOCKS : "shadowsocks",
    //文章类型
    ARTICLE_TYPE_HOME      : "首页",
    ARTICLE_TYPE_DIALOG    : "弹窗",
    ARTICLE_TYPE_NOTICE    : "公告",
  },
};