export default {
  adminNode:{
    basicParameters:"基础参数",
    otherParameters:"其他参数",
    //
    addNode:"新增节点",
    modifyNode:"修改节点",
    sortNode:"排序节点",
    sharedNode:"共享节点",
    parse:"解析",
    placeholder1:"支持订阅、单个节点、多个节点、base64编码解析（vmess，vless，trojan，shadowsocks, hy2）",
    serverStatus:"服务器状态",
    NodeInfo: {
      node_type_transfer:"中转",
      node_type_direct:"直连",
      index:"序号",
      created_at: "创建时间",
      updated_at: "更新时间",
      id: "id",
      node_speed_limit: "节点限速", //节点限速/Mbps
      traffic_rate: "倍率",    //倍率
      node_type: "节点类型",       //类型 vless vmess shadowsocks
      //基础参数
      remarks: "别名",//别名
      address: "地址",
      port: "端口",
      node_order: "节点排序",//节点排序
      enabled: "状态",  //是否为激活节点
      protocol:"节点协议",
      //中转参数
      enable_transfer: "启用中转",//是否启用中转
      transfer_address: "中转地址",//中转ip
      transfer_port: "中转端口",   //中转port
      transfer_node_id: "绑定的父节点", //中转绑定的节点ID
      // 共享节点额外需要的参数
      is_shared_node:"是否共享节点",
      uuid:"uuid",
      //流量
      total_up: "上行流量",
      total_down: "下行流量",
      //vmess参数
      v: "v",
      scy: "scy",//加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
      aid: "aid",//额外ID
      //vless参数
      flow: "flow",//流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
      encryption: "encryption",//加密方式 none

      network: "network",//传输协议 tcp,kcp,ws,h2,quic,grpc
      type: "type",   //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
      host: "host",   //伪装域名
      path: "path",   //path
      mode: "mode",   //grpc传输模式 gun，multi
      service_name: "service_name",

      allowInsecure: "allowInsecure",//tls 跳过证书验证
      security: "security",//传输层安全类型 none,tls,reality
      sni: "sni",
      fp: "fp",
      alpn: "alpn",
      dest: "dest",
      private_key: "private_key",
      pbk: "pbk",
      sid: "sid",
      spx: "spx",
      //关联参数
      goods:"关联商品",//关联商品
      access: "访问控制",
    },
    //
    access:"访问控制",
    addAccess:"新增访问控制",
    modifyAccess:"修改访问控制",
    Access:{
      index:"序号",
      created_at:"创建时间",
      updated_at:"更新时间",
      id:"id",
      name:"名称",
      route:"禁止路由规则",
    },
  },
}