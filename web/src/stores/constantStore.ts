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
    NODE_TYPE_TRANSFER    : "transfer",
    NODE_TYPE_VMESS       : "vmess",
    NODE_TYPE_VLESS       : "vless",
    NODE_TYPE_TROJAN      : "trojan",
    NODE_TYPE_HYSTERIA    : "hysteria",
    NODE_TYPE_SHADOWSOCKS : "shadowsocks",
    //文章类型
    ARTICLE_TYPE_HOME      : "home",
    ARTICLE_TYPE_DIALOG    : "dialog",
    ARTICLE_TYPE_NOTICE    : "notice",
    ARTICLE_TYPE_KNOWLEDGE : "knowledge",


  }),
  actions: {

  }
})