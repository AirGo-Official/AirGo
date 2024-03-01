import {defineStore, storeToRefs} from "pinia";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


export const useShopStore = defineStore("shopStore", {
    state: () => ({
        //当前支付商品
        currentGoods: {
            id: 0,
            created_at: "",
            updated_at: "",
            good_order: 0,
            is_show: false,
            des: '',
            subject: "",
            price: "",
            goods_type: '',
            deliver_type: '',
            deliver_text: '',
            node_connector: 0,
            total_bandwidth: 0,
            recharge_amount: '',
            nodes: [] as NodeInfo[],
        } as Goods,
        //当前商品订单
        currentOrder: {
            order_type:'', //订单类型:New=新购入;Renew=续费"
            trade_status: '', //交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"
            out_trade_no: '', //商户订单号
            total_amount: '', //订单金额
            buyer_pay_amount: '', //付款金额
            coupon_amount: '', //折扣码这块金额
            balance_amount: '', //余额折扣金额
            duration:1,// 默认订购时长
            // 关联用户
            user_id: 0,
            user_name: '',
            user: {},
            // 商品参数
            goods_id: 0,
            goods_type: '',  //类型
            deliver_type: '', //发货类型
            deliver_text: '', //发货内容
            subject: '',
            price: '',

            //支付参数
            pay_id: 0,   //支付方式id
            pay_type: '', //支付方式，alipay,epay
            coupon_id: 0,
            coupon_name: '',
            pay_info: {
                alipay_info:{qr_code: ''},
                epay_info:{
                    epay_api_url: '',
                    epay_pre_create_pay: {
                        pid: 0,
                        type: '',
                        out_trade_no: '',
                        notify_url: '',
                        return_url: '',
                        name: '',
                        money: '',
                        clientip: '',
                        device: '',
                        param: '',
                        sign: '',
                        sign_type: '',
    },
},
} as PreCreatePayToFrontend, //支付信息，epay，alipay等
            trade_no: '',
            buyer_logon_id: '',
        } as Order,
        //商品列表
        goodsList: [] as Goods[],
        orderList: {
            data: [] as Order[],
            total: 0,
        },
        payList:[] as PayInfo[]
    }),
    actions: {
        //加载时获取全部已启用商品
        async getAllEnabledGoods(params:object) {
            const res = await request(apiStoreData.userApi.value.getEnabledGoodsList,params)
            this.goodsList = res.data
        },
        //获取用户订单
        async getOrderList(params?: object) {
            const res = await request(apiStoreData.userApi.value.getOrderList, params)
            this.orderList = res.data
        },
        //获取订单详情(下单时）
        async getOrderInfo(params: object) {
            return  request(apiStoreData.userApi.value.getOrderInfo, params)
        },
        //获取启用的支付列表
        async getEnabledPayList() {
            const res = await request(apiStoreData.userApi.value.getEnabledPayList)
            this.payList = res.data
        },
    }
})