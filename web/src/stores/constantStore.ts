import {defineStore} from 'pinia';
export const useConstantStore = defineStore('constantStore', {
  state: () => ({
    // 商品
    GOODS_TYPE_GENERAL   : "general",
    GOODS_TYPE_SUBSCRIBE : "subscribe",
    GOODS_TYPE_RECHARGE  : "recharge",
    // 订单类型
    ORDER_TYPE_NEW:"New",
    ORDER_TYPE_RENEW:"Renew",
    ORDER_TYPE_DESTROYED:"Destroyed",
    //订单默认超时时间
    CACHE_SUBMIT_ORDER_TIMEOUT  : 10,    //Minute
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
    EMAIL_TYPE_TEST : "EMAIL_TYPE_TEST",
    //节点类型
    NODE_TYPE_NORMAL   : "normal",
    NODE_TYPE_TRANSFER : "transfer",
    NODE_TYPE_SHARED   : "shared",
    //节点协议类型
    NODE_PROTOCOL_VMESS       : "vmess",
    NODE_PROTOCOL_VLESS       : "vless",
    NODE_PROTOCOL_TROJAN      : "trojan",
    NODE_PROTOCOL_HYSTERIA    : "hysteria2",
    NODE_PROTOCOL_SHADOWSOCKS : "shadowsocks",

    // 支付类型
    PAY_TYPE_ALIPAY  : "alipay",  // 支付宝alipay
    PAY_TYPE_EPAY    : "epay",    //易支付
    PAY_TYPE_BALANCE : "balance", //余额支付

    //文章类型
    ARTICLE_TYPE_HOME      : "home",
    ARTICLE_TYPE_DIALOG    : "dialog",
    ARTICLE_TYPE_NOTICE    : "notice",
    ARTICLE_TYPE_KNOWLEDGE : "knowledge",

    //finance
    BALANCE_STATEMENT_TITLE_RECHARGE    : "Recharge",
    BALANCE_STATEMENT_TITLE_EXPENDITURE : "Expenditure",
    BALANCE_STATEMENT_TITLE_WITHDRAW    : "Withdraw",
    BALANCE_STATEMENT_TYPE_PLUS         : "Plus",
    BALANCE_STATEMENT_TYPE_REDUCE       : "Reduce",

    //业务代码
    RESPONSE_ERROR   : 1,  //code=1，不能正常获取请求数据
    RESPONSE_SUCCESS : 0,  //code=0，能正常获取请求数据
    RESPONSE_WARNING : 10, //code=10，能正常获取请求数据，但有重要message 需要显式提醒
    TOKENERROR : 40101, //token过期
    LIMITERROR : 40201, //限流

  }),
  actions: {

  }
})