import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const usePayStore = defineStore('payStore', {
    state: () => ({
        //当前编辑的pay
        payInfo: {
            id: 0,
            name: '',
            pay_type: '',
            pay_logo_url: '',
            status: false,
            alipay: {
                alipay_app_id: '',
                alipay_return_url: '',
                alipay_app_private_key: '',
                alipay_ali_public_key: '',
                alipay_encrypt_key: '',
            } as Alipay,
            epay: {
                epay_pid: 0,
                epay_key: '',
                epay_api_url: '',
                epay_return_url: '',
                epay_notify_url: '',
                epay_type: '',
            } as Epay,
        } as PayInfo,
        //pay list
        payList: [] as PayInfo[],
    }),
    actions: {
        //获取支付列表
        async getPayList() {
            const res = await request(apiStoreData.api.value.pay_getPayList)
            this.payList = res.data
        },
        //获取启用的支付列表
        async getEnabledPayList() {
            const res = await request(apiStoreData.api.value.pay_getEnabledPayList)
            this.payList = res.data
        },
        //清空payInfo
        wipePayInfo() {
            this.payInfo = {
                id: 0,
                name: '',
                pay_type: '',
                pay_logo_url: '',
                status: false,
                alipay: {
                    alipay_app_id: '',
                    alipay_return_url: '',
                    alipay_app_private_key: '',
                    alipay_ali_public_key: '',
                    alipay_encrypt_key: '',
                } as Alipay,
                epay: {
                    epay_pid: 0,
                    epay_key: '',
                    epay_api_url: '',
                    epay_return_url: '',
                    epay_notify_url: '',
                    epay_type: '',
                } as Epay,
            } as PayInfo
        },
    }
})