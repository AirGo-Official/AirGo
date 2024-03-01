export default {
  adminNode:{
   basicParameters:"Basic parameters",
   otherParameters:"Other parameters",
   //
   addNode:"Add node",
   modifyNode:"Modify node",
   sortNode:"Sort node",
   sharedNode:"Shared node",
    parse:"Parse",
    placeholder1:"Supports subscription, single node, multiple nodes, and base64 encoding resolution (vmess, vless, trojan, shadowlocks, hy2)",
    serverStatus:"Server Status",
   NodeInfo: {
      node_type_transfer:"Transfer",
      node_type_direct:"Direct",
      index:"index",
      created_at: "created_at",
      updated_at: "updated_at",
      id: "id",
      node_speed_limit: "node_speed_limit", //节点限速/Mbps
      traffic_rate: "traffic_rate",    //倍率
      node_type: "node_type",       //类型 vless vmess shadowsocks
      //基础参数
      remarks: "remarks",//别名
      address: "address",
      port: "port",
      node_order: "node order",//节点排序
      enabled: "status",  //是否为激活节点
     protocol:"protocol",
      //中转参数
      enable_transfer: "enable transfer",//是否启用中转
      transfer_address: "transfer address",//中转ip
      transfer_port: "transfer port",   //中转port
      transfer_node_id: "transfer parent node", //中转绑定的节点ID
     // 共享节点额外需要的参数
     is_shared_node:"是否共享节点",
     uuid:"uuid",
      //流量
      total_up: "total up",
      total_down: "total down",

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
      goods:"goods",//关联商品
      access: "access",
    },
   //
   access:"Access",
   addAccess:"Add Access",
   modifyAccess:"Modify Access",
   Access:{
    index:"index",
    created_at:"created at",
    updated_at:"updated at",
    id:"id",
    name:"name",
    route:"route",
   },
  },
}