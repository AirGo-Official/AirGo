import {defineStore, storeToRefs} from "pinia";

import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


export const useShopStore = defineStore("shopStore", {
    state: () => ({
        //商品管理页面数据
        goodsManageData: {
            //当前编辑商品
            currentGoods: {
                subject: "新套餐",
                total_amount: "0.00",
                total_bandwidth: 0,
                expiration_date: 0,
                // checked_nodes: [0], //套餐编辑时选中的节点
                // nodes: [],
                des: '<h3 style="color:#00BFFF">究竟什么样的终点，才配得上这一路的颠沛流离---管泽元</h3>\n<h3 style="color:#DDA0DD">世界聚焦于你---管泽元</h3>',
            } as Goods,
        },
        //全部商品
        goodsList: [] as Goods[],
        //编辑折扣码选中的商品

        //商店页面参数
        shopData: {
            //当前支付商品
            currentGoods: {
                id: 0, //不能覆盖
                created_at: "",
                updated_at: "",
                good_order: 0,
                status: false,
                des: '',
                subject: "",
                total_amount: "",
                product_code: "",
                total_bandwidth: 0,
                expiration_date: 0,
                checked_nodes: [0],
                nodes: [],

            } as Goods,
            //当前商品订单
            currentOrder: {
                id: 0,
                out_trade_no: '',
                goods_id: 0,
                subject: '',
                price: '',
                pay_type: '',

                coupon: 0,
                coupon_name: '',
                coupon_amount: '0',
                deduction_amount: '0',

                trade_no: '',
                buyer_logon_id: '',
                trade_status: '',
                total_amount: '',
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
                },
            } as Order,
        }
    }),
    actions: {
        //加载时获取全部已启用商品
        async getAllEnabledGoods() {
            const res = await request(apiStoreData.api.value.shop_getAllEnabledGoods)
            this.goodsList = res.data
            ElMessage.success(res.msg)
        },
        //获取全部订阅商品
        async getAllGoods() {
            const res = await request(apiStoreData.api.value.shop_getAllGoods)
            this.goodsList = res.data
        },
        //添加商品
        async newGoods(data?: object) {
            const res = await request(apiStoreData.api.value.shop_newGoods,data)
            ElMessage.success(res.msg)
        },
        //修改商品
        async updateGoods(data?: object) {
            const res = await request(apiStoreData.api.value.shop_updateGoods,data)
            ElMessage.success(res.msg)
        },
        //删除商品
        async deleteGoods(data?: object) {
            const res = await request(apiStoreData.api.value.shop_deleteGoods,data)
            ElMessage.success(res.msg)
        },
    }
})