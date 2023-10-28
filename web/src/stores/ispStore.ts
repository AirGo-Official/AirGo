import {defineStore, storeToRefs} from "pinia";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

export const useISPStore = defineStore("ispStore", {
    state: () => ({
        isp: {
            user_id: 0,
            isp_type: '',
            status: false,
            mobile: '',
            unicom_config: {
                version: '',
                app_id: '',
                cookie: '',
                unicomMobile: '',
                password: '',
            },
            telecom_config: {
                phoneNum: '',
                telecomPassword: '',
                timestamp: '',
                loginAuthCipherAsymmertric: '',
                deviceUid: '',
                telecomToken: '',
            },

        } as Isp,
        isCountDown: false,
        countDownTime: 60,
    }),
    actions: {
        async sendCode(params?: object) {
            const res = await request(apiStoreData.api.value.isp_sendCode, params)
            ElMessage.success(res.msg)
        },
        async ispLogin(params?: object) {
            const res = await request(apiStoreData.api.value.isp_ispLogin, params)
            if (res.msg === '获取成功') {
                this.isp = res.data
                return
            }
            this.getMonitorByUserID()
        },
        async getMonitorByUserID(params?: object) {
            const res = await await request(apiStoreData.api.value.isp_getMonitorByUserID, params)
            ElMessage.success(res.msg)
            this.isp = res.data
        },
    }
})