import {defineStore, storeToRefs} from "pinia";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { useConstantStore } from "/@/stores/constantStore";
const constantStore = useConstantStore()
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

export const useAdminNodeStore = defineStore("adminNodeStore", {
    state: () => ({
        //节点列表
        nodeList: {
            total: 0,
            data: [] as NodeInfo[],
        },
        //弹窗页数据
        dialogData: {
            checkedAccessIDs:[] as number[],
            vlessInfo: {
                node_speed_limit: 0, //节点限速/Mbps
                traffic_rate: 1,    //倍率
                node_type: 'vless',
                uuid: '',
                //基础参数
                remarks: '',//别名
                address: '',
                port: 80,
                node_order: 0,//节点排序
                enabled: true,  //是否为激活节点
                //中转参数
                enable_transfer: false,//是否启用中转
                transfer_address: '',//中转ip
                transfer_port: 0,   //中转port
                total_up: 0,
                total_down: 0,
                goods: [] as Goods[],//多对多,关联商品
                v: '2',
                scy: 'none',
                aid: 0,//额外ID
                flow: '',//流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
                encryption: 'none',//加密方式
                network: 'ws',//传输协议
                type: '',   //伪装类型
                host: '',   //伪装域名
                path: '/',   //path
                mode: 'multi',   //grpc传输模式 gun，multi
                service_name: '',
                security: 'none',//传输层安全类型 none,tls,reality
                sni: '',
                fp: 'chrome',
                alpn: '',
                allowInsecure: true,//tls 跳过证书验证
                dest: '',
                private_key: '',
                pbk: '',
                sid: '',
                spx: '',
                access:[] as Access[],
            } as NodeInfo,
            vmessInfo: {
                node_speed_limit: 0, //节点限速/Mbps
                traffic_rate: 1,    //倍率
                node_type: 'vmess',
                uuid: '',
                //基础参数
                remarks: '',//别名
                address: '',
                port: 80,
                node_order: 0,//节点排序
                enabled: true,  //是否为激活节点
                //中转参数
                enable_transfer: false,//是否启用中转
                transfer_address: '',//中转ip
                transfer_port: 0,   //中转port
                //
                total_up: 0,
                total_down: 0,
                goods: [] as Goods[],//多对多,关联商品
                //vmess参数
                v: '2',
                scy: 'auto',//加密方式
                aid: 0,//额外ID
                //vless参数
                flow: '',//流控
                encryption: '',//加密方式 none
                network: '',//传输协议
                type: '',   //伪装类型
                host: '',   //伪装域名
                path: '/',   //path
                mode: 'multi',   //grpc传输模式 gun，multi
                service_name: '',
                security: 'none',//传输层安全类型 none,tls,reality
                sni: '',
                fp: 'chrome',
                alpn: '',
                allowInsecure: true,//tls 跳过证书验证
                dest: '',
                private_key: '',
                pbk: '',
                sid: '',
                spx: '',
                access:[] as Access[],
            } as NodeInfo,
            shadowsocksInfo: {
                node_speed_limit: 0, //节点限速/Mbps
                traffic_rate: 1,    //倍率
                node_type: 'shadowsocks',
                uuid: '',
                //基础参数
                remarks: '',//别名
                address: '',
                port: 80,
                node_order: 0,//节点排序
                enabled: true,  //是否为激活节点
                //中转参数
                enable_transfer: false,//是否启用中转
                transfer_address: '',//中转ip
                transfer_port: 0,   //中转port
                //
                total_up: 0,
                total_down: 0,
                goods: [] as Goods[],//多对多,关联商品
                //vmess参数
                v: '2',
                scy: 'none',//加密方式
                aid: 0,//额外ID
                //vless参数
                flow: '',//流控
                encryption: '',//加密方式 none
                network: '',//传输协议
                type: '',   //伪装类型
                host: '',   //伪装域名
                path: '/',   //path
                mode: 'multi',   //grpc传输模式 gun，multi
                service_name: '',
                security: 'none',//传输层安全类型 none,tls,reality
                sni: '',
                fp: 'chrome',
                alpn: '',
                allowInsecure: true,//tls 跳过证书验证
                dest: '',
                private_key: '',
                pbk: '',
                sid: '',
                spx: '',
                access:[] as Access[],
            } as NodeInfo,
            hysteriaInfo: {
                node_speed_limit: 0, //节点限速/Mbps
                traffic_rate: 1,    //倍率
                node_type: 'hysteria',
                uuid: '',
                //基础参数
                remarks: '',//别名
                address: '',
                port: 80,
                node_order: 0,//节点排序
                enabled: true,  //是否为激活节点
                //中转参数
                enable_transfer: false,//是否启用中转
                transfer_address: '',//中转ip
                transfer_port: 0,   //中转port
                //
                total_up: 0,
                total_down: 0,
                goods: [] as Goods[],//多对多,关联商品
                //vmess参数
                v: '2',
                scy: 'none',//加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
                aid: 0,//额外ID
                //vless参数
                flow: '',//流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
                encryption: '',//加密方式 none
                network: '',//传输协议 tcp,kcp,ws,h2,quic,grpc
                type: '',   //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
                host: '',   //伪装域名
                path: '/',   //path
                mode: 'multi',   //grpc传输模式 gun，multi
                service_name: '',
                security: 'none',//传输层安全类型 none,tls,reality
                sni: '',
                fp: 'chrome',
                alpn: '',
                allowInsecure: true,//tls 跳过证书验证
                dest: '',
                private_key: '',
                pbk: '',
                sid: '',
                spx: '',
                access:[] as Access[],
            } as NodeInfo,
            transferInfo: {
                remarks: '',
                enabled: true,
                node_type: 'transfer',
                enable_transfer: true,
                transfer_address: '',
                transfer_port: 0,
                transfer_node_id: 0,
            } as NodeInfo,
        },
        //节点状态页面数据
        serverStatusData: {
            type: 0,
            data: [] as ServerStatusInfo[],
        },
        //共享节点页面数据
        nodeSharedData: {
            newNodeSharedUrl: {
                url: '',
            },
            nodeList: {
                total:0,
                data:[] as NodeSharedInfo[],
            },
        },
        accessList: {
            total: 0,
            data: [] as Access[]
        },
        currentAccess: {} as Access,
    }),
    actions: {
        //获取全部节点
        async getAllNode() {
            const res = await request(apiStoreData.adminApi.value.getAllNode)
            this.nodeList = res.data
        },
        //获取全部节点 with Traffic,分页
        async getNodeWithTraffic(params?: object) {
            const res = await request(apiStoreData.adminApi.value.getNodeListWithTraffic, params)
            this.nodeList = res.data
            //处理节点的access

        },
        //更新节点
        async updateNode(params: NodeInfo) {
            return  request(apiStoreData.adminApi.value.updateNode, this.accessHandlerWhenSubmit(params))
        },
        //删除节点
        async deleteNode(params: NodeInfo) {
            return  request(apiStoreData.adminApi.value.deleteNode, params)
        },
        //新建节点
        async newNode(params: NodeInfo) {
            return  request(apiStoreData.adminApi.value.newNode, this.accessHandlerWhenSubmit(params))
        },
        //新建共享节点
        async newNodeShared(params: object) {
            return  request(apiStoreData.adminApi.value.newNodeShared, params)
        },
        //获取共享节点列表
        async getNodeSharedList() {
            const res = await request(apiStoreData.adminApi.value.getNodeSharedList)
            this.nodeSharedData.nodeList = res.data
        },
        //删除共享节点
        async deleteNodeShared(params: NodeInfo) {
            return  request(apiStoreData.adminApi.value.deleteNodeShared, params)
        },
        //根据节点类型返回节点对象
        returnNodeInfo(nodeType: string) {
            let n = {} as NodeInfo
            switch (nodeType) {
                case constantStore.NODE_TYPE_VLESS:
                    n= this.dialogData.vlessInfo
                    break
                case constantStore.NODE_TYPE_VMESS:
                    n = this.dialogData.vmessInfo
                    break
                case constantStore.NODE_TYPE_SHADOWSOCKS:
                    n = this.dialogData.shadowsocksInfo
                    break
                case constantStore.NODE_TYPE_HYSTERIA:
                    n = this.dialogData.hysteriaInfo
                    break
                case constantStore.NODE_TYPE_TRANSFER:
                    n = this.dialogData.transferInfo
                    break
            }
            return n
        },
        accessHandler(node:NodeInfo){
            this.dialogData.checkedAccessIDs = []
            node.access.forEach((value: Access, index: number, array: Access[])=>{
                this.dialogData.checkedAccessIDs.push(value.id)
            })
        },
        accessHandlerWhenSubmit(node:NodeInfo){
            node.access = [] as Access[]
            this.dialogData.checkedAccessIDs.forEach((value: number, index: number, array: number[])=>{
                node.access.push({id:value} as Access)
            })
            return node
        },

        async getAccessList(params:QueryParams) {
            const res = await request(apiStoreData.adminApi.value.getAccessList,params)
            this.accessList = res.data
        },
        async newAccess(params: Access) {
            return  request(apiStoreData.adminApi.value.newAccess, params)
        },
        async deleteAccess(params: Access) {
            return  request(apiStoreData.adminApi.value.deleteAccess, params)
        },
        async updateAccess(params: Access) {
            return  request(apiStoreData.adminApi.value.updateAccess, params)
        }
    }
})