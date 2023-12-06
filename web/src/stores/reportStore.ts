import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useReportStore = defineStore("reportStore", {
    state: () => ({
        //所有需要显示的数据库表
        allDbTables: [
            {en_name: "user", cn_name: "用户"},
            {en_name: "orders", cn_name: "订单"},
            {en_name: "goods", cn_name: "商品"},
            {en_name: "node", cn_name: "节点"},
            {en_name: "pay", cn_name: "支付"},
            {en_name: "article", cn_name: "文章"},
            {en_name: "coupon", cn_name: "折扣码"},
            {en_name: "role", cn_name: "角色"},
            {en_name: "access", cn_name: "访问控制"},
            {en_name: "gallery", cn_name: "图库列表"},
        ],
        //选中的数据库，库表，用来请求获取数据库的数据表的所有字段名,类型值
        checkedDbInfo: {
            db_name: '',
            table_name: '',
        },
        //高级查询的条件参数
        reportParams: {
            table_name: '',
            field_params_list: [] as FieldParams[],
            pagination: {
                page_num: 1,
                page_size: 30,
                order_by: '',
            } as Pagination,//分页参数
        },
        //字段信息
        fieldData: {
            field_list: [],
            field_chinese_name_list: {} as { [key: string]: any; },
            field_type_list: {} as { [key: string]: any; },
        },
        //保存数据
        reportData: {
            total: 0,
            data: [],
        }

    }),
    actions: {
        // 获取字段名,类型值
        async getColumn() {
            const res = await request(apiStoreData.api.value.report_getColumn, this.checkedDbInfo)
            this.fieldData = res.data
        },
        //查询
        async getReport() {
            const res = await request(apiStoreData.api.value.report_reportSubmit, this.reportParams)
            this.reportData = res.data
        },
    }
})