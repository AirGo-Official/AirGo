import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


export const usePublicStore = defineStore('publicStore', {
    state: () => ({
        base64CaptchaData: {
            id: '',
            b64s: '',
        } as Base64CaptchaInfo,
        publicSetting:{
            enable_register: true,          // 是否开启注册
            enable_email_code: false,       //是否开启注册邮箱验证码
            enable_login_email_code: false, //是否开启登录邮箱验证码
            acceptable_email_suffixes: '', //可接受的邮箱后缀
            rebate_rate: 0,                  //佣金率
            backend_url: '',                 //
            enabled_clock_in: true,           //是否开启打卡
        } ,
        acceptable_email_suffixes_arr: [] as string[],    //可接受的邮箱后缀数组

    }),
    actions: {
        //获取base64Captcha
        async getBase64Captcha() {
            const res = await request(apiStoreData.publicApi.value.getBase64Captcha)
            this.base64CaptchaData = res.data
        },
        //获取公共系统设置
        async getPublicSetting(){
            const res = await request(apiStoreData.publicApi.value.getPublicSetting)
            this.publicSetting = res.data
            this.acceptable_email_suffixes_arr = this.publicSetting.acceptable_email_suffixes.split("\n")
        },
        //发送验证码
        async sendEmailCode(params:EmailRequest) {
            return request(apiStore.publicApi.getEmailCode, params)
        },

    }

})