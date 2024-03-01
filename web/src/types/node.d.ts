declare interface NodeInfo {
    created_at: string;
    updated_at: string;
    id: number;
    node_speed_limit: number; //节点限速/Mbps
    traffic_rate: number;    //倍率
    node_type: string;       //类型 vless vmess shadowsocks
    //基础参数
    remarks: string;//别名
    address: string;
    port: number;
    node_order: number;//节点排序
    enabled: boolean;  //是否为激活节点（是否显示）
    //中转参数
    transfer_address: string;//中转ip
    transfer_port: number;   //中转port
    transfer_node_id: number; //中转绑定的节点ID
    // 共享节点额外需要的参数
    uuid:string
    //流量
    total_up: number;
    total_down: number;

    //节点协议类型
    protocol:string

    //vmess参数
    v: string;
    scy: string;//加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
    aid: number;//额外ID
    //vless参数
    flow: string;//流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
    encryption: string;//加密方式 none

    network: string;//传输协议 tcp,kcp,ws,h2,quic,grpc
    type: string;   //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
    host: string;   //伪装域名
    path: string;   //path
    mode: string;   //grpc传输模式 gun，multi
    service_name: string;

    allowInsecure: boolean;//tls 跳过证书验证
    security: string;//传输层安全类型 none,tls,reality
    sni: string;
    fp: string;
    alpn: string;
    dest: string;
    private_key: string;
    pbk: string;
    sid: string;
    spx: string;
    //关联参数
    goods: goods [];//多对多,关联商品
    access: Access[];
}

declare interface ServerStatusInfo {
    id: number;
    name: string;
    status: boolean;
    last_time: string;
    user_amount: number;
    traffic_rate: number;
    u: number;
    d: number;
    cpu: number;
    mem: number;
    disk: number;
}

declare interface RealityItem {
    dest: string
    sni: string
}
