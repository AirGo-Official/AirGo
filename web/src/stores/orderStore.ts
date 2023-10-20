import {defineStore, storeToRefs} from "pinia";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useOrderStore = defineStore("orderStore", {
    state: () => ({
        //订单管理页面数据
        orderManageData: {
            allOrders: {
                order_list: [] as Order[],
                total: 0,
            },
        },
        //个人订单数据
        orderPersonal: {
            allOrders: {
                order_list: [] as Order[],
                total: 0,
            },
        },
    }),
    actions: {
        //获取订单详情(下单时）
        async getOrderInfo(params: object) {
            const res = await request(apiStoreData.api.value.order_getOrderInfo, params)
            return res
        },
        //获取全部订单
        async getAllOrder(params?: object) {
            const res = await request(apiStoreData.api.value.order_getAllOrder, params)
            this.orderManageData.allOrders = res.data
            ElMessage.success(res.msg)
        },
        //获取用户最近10次订单
        async getOrder(params?: object) {
            const res = await request(apiStoreData.api.value.order_getOrderByUserID, params)
            this.orderPersonal.allOrders.order_list = res.data
            ElMessage.success(res.msg)
        },
        //完成未支付订单
        async completedOrder(params?: object) {
            const res = await request(apiStoreData.api.value.order_completedOrder, params)
            ElMessage.success(res.msg)
        },
    }
})