export default {
  common: {
    tip:"Tip",
    message_confirm_delete:"This operation will permanently delete the data. Do you want to continue?",
    button_cancel:"Cancel",
    button_confirm:"Confirm",

    operate:"Operate",
    delete:"Delete",
    modify:"Modify",
    reset:"Reset",
    query:"Query",
    advanced_query:"Advanced query",
    add:"Add",
    sort:"Sort",
    reply:"Reply",
    close:"Close",
    back:"back",
    send:"Send",
    details:"details",
    purchase:"purchase",

    yes:"Yes",
    no:"No",
    enable:"Enable",
    disable:"Disable",
    display:"display",
    hide:"hide",

    startDate:"Start date",
    endDate:"End date",
    to:"to",
    lastWeek:"Last week",
    lastMonth:"Last month",
    lastThreeMonths:"Last 3 months",
    month:"month",

    noData:"No data",

    success:"Success",
    warning:"Warning",
    error:"Error",
  },
  constant:{
    // 商品
    GOODS_TYPE_GENERAL   : "general",
    GOODS_TYPE_SUBSCRIBE : "subscribe",
    GOODS_TYPE_RECHARGE  : "recharge",
    // 订单类型
    ORDER_TYPE_NEW:"New",
    ORDER_TYPE_RENEW:"Renew",
    ORDER_TYPE_DESTROYED:"Destroyed",
    // 订单状态
    ORDER_STATUS_CREATED        : "CREATED",
    ORDER_STATUS_WAIT_BUYER_PAY : "WAIT_BUYER_PAY",
    ORDER_STATUS_TRADE_CLOSED   : "TRADE_CLOSED",
    ORDER_STATUS_TRADE_SUCCESS  : "TRADE_SUCCESS",
    ORDER_STATUS_TRADE_FINISHED : "TRADE_FINISHED",
    ORDER_STATUS_COMPLETED      : "COMPLETED ",
    ORDER_STATUS_UNKNOWN_STATE:"UNKNOWN_STATE",
    // 发货
    DELIVER_TYPE_NONE   : "none",
    DELIVER_TYPE_AUTO   : "auto",
    DELIVER_TYPE_MANUAL : "manual",
    //工单状态
    TICKET_PROCESSING : "TICKET_PROCESSING",//工单进行中
    TICKET_CLOSED     :"TICKET_CLOSED",
    //邮件类型
    EMAIL_TYPE_USER_REGISTER : "EMAIL_TYPE_USER_REGISTER",
    EMAIL_TYPE_USER_RESETPWD : "EMAIL_TYPE_USER_RESETPWD",
    //节点类型
    NODE_TYPE_NORMAL    : "normal",
    NODE_TYPE_SHARED    : "shared",
    NODE_TYPE_TRANSFER    : "transfer",
    //节点协议类型
    NODE_PROTOCOL_VMESS       : "vmess",
    NODE_PROTOCOL_VLESS       : "vless",
    NODE_PROTOCOL_TROJAN      : "trojan",
    NODE_PROTOCOL_HYSTERIA    : "hysteria",
    NODE_PROTOCOL_SHADOWSOCKS : "shadowsocks",
    //文章类型
    ARTICLE_TYPE_HOME      : "home",
    ARTICLE_TYPE_DIALOG    : "dialog",
    ARTICLE_TYPE_NOTICE    : "notice",
  },
};