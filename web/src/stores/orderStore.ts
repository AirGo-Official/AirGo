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
            fieldParams:{
                table_name: 'orders',
                field_params_list: [
                    {field:"out_trade_no",field_chinese_name:"",field_type:"",condition:"like",condition_value:"",operator:""}
                ] as FieldParams[],
                pagination: {
                    page_num: 1,
                    page_size: 30,
                    order_by: 'id DESC',
                } as Pagination,//分页参数
            },
            currentOrder: {
                created_at: '',
                updated_at: '',
                id: 0,
                user_id: 0,
                user_name: '',
                // user: any;

                out_trade_no: '',
                goods_id: 0,
                goods_type: '',  //类型
                deliver_type: '',//发货类型
                deliver_text: '',//发货内容
                subject: '',
                price: '',
                pay_id: 0,   //支付方式id
                pay_type: '', //支付方式，alipay,epay
                coupon_id: 0,    //
                coupon_name: '',
                coupon_amount: '',
                deduction_amount:'',
                remain_amount: '',

                // pay_info: {} as PreCreatePayToFrontend, //支付信息，epay，alipay等
                trade_no: '',
                buyer_logon_id: '',
                trade_status: '',
                total_amount: '',
                receipt_amount: '',
                buyer_pay_amount: '',
            } as Order,
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
            this.orderManageData.allOrders.order_list = res.data.data
            this.orderManageData.allOrders.total = res.data.total
        },
        //获取用户订单
        async getOrder(params?: object) {
            const res = await request(apiStoreData.api.value.order_getOrderByUserID, params)
            this.orderPersonal.allOrders.order_list = res.data.data
            this.orderPersonal.allOrders.total = res.data.total
        },
        //完成未支付订单
        async completedOrder(params?: object) {
            const res = await request(apiStoreData.api.value.order_completedOrder, params)
        },
        //更新用户订单
        async updateUserOrder(params:object){
            const res = await request(apiStoreData.api.value.order_updateUserOrder, params)
        }
    }
})