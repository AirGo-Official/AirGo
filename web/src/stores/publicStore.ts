import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { Session } from "/@/utils/storage";

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
            enable_base64_captcha:true, //是否开启注册图片验证码
            enable_email_code: false,       //是否开启注册邮箱验证码
            enable_login_email_code: false, //是否开启登录邮箱验证码
            acceptable_email_suffixes: '', //可接受的邮箱后缀
            backend_url: '',                 //
            enabled_clock_in: true,           //是否开启打卡
            commission_rate:0,                //佣金率, 范围 0~1, 佣金 = 订单金额 * 佣金率 ( 100.50 * 0.50 )
            withdraw_threshold:0,            //提取到余额的阈值
            enable_lottery:false,
            jackpot:[] as JackpotItem[],
            sub_name:'',
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
            //尝试从session中获取
            if (Session.get('publicSetting')){
                this.publicSetting =  Session.get('publicSetting')
                this.acceptable_email_suffixes_arr = this.publicSetting.acceptable_email_suffixes.split("\n")
            } else {
                const res = await request(apiStoreData.publicApi.value.getPublicSetting)
                this.publicSetting = res.data
                this.acceptable_email_suffixes_arr = this.publicSetting.acceptable_email_suffixes.split("\n")
                // Session.set('publicSetting',this.publicSetting) //关闭缓存
            }
        },
        //发送验证码
        async sendEmailCode(params:EmailRequest) {
            return request(apiStore.publicApi.getEmailCode, params)
        },
    }
})