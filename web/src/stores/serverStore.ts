import {defineStore, storeToRefs} from "pinia";

import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

export const useServerStore = defineStore("serverStore", {
    state: () => ({
        serverConfig: {
            created_at: '',
            updated_at: '',
            id: 0,
            security:{
                jwt: {
                    signing_key: '',
                    expires_time: '',
                    buffer_time: '',
                    issuer: '',
                },
                captcha: {
                    key_long: 0,
                    img_width: 0,
                    img_height: 0,
                    open_captcha: 0,
                    open_captcha_time_out: 0,
                },
                rate_limit_params: {
                    ip_role_param: 0,
                    visit_param: 0,
                },
            } as Security,
            subscribe: {
                enable_register: true,
                enable_email_code: false,
                enable_login_email_code: false,
                acceptable_email_suffixes:'',
                is_multipoint: true,
                sub_name: '',
                backend_url: '',
                api_prefix: '',
                tek: '',
                default_goods: '',
                enabled_rebate: true,    //是否开启返利
                rebate_rate: 0,          //返利率
                enabled_deduction: true, //是否开启旧套餐抵扣
                deduction_threshold: 0,  //旧套餐抵扣阈值,大于该值则抵扣
                enabled_clock_in: true,
                clock_in_min_traffic: 100,
                clock_in_max_traffic: 1000,
            },

            email: {
                email_from: '',
                email_secret: '',
                email_host: '',
                email_port: 0,
                email_is_ssl: true,
                email_nickname: '',
                email_subject: '',
                email_content: '',
            },

        } as Server,
        publicServerConfig: {
            enable_register: true,          // 是否开启注册
            enable_email_code: false,       //是否开启注册邮箱验证码
            enable_login_email_code: false, //是否开启登录邮箱验证码
            acceptable_email_suffixes: '', //可接受的邮箱后缀
            rebate_rate: 0,                  //佣金率
            backend_url: '',                 //
            enabled_clock_in:true,           //是否开启打卡
        },
        acceptable_email_suffixes_arr:[] as string[],    //可接受的邮箱后缀数组

    }),
    actions: {
        //获取系统设置
        async getServerConfig() {
            const apiStore = useApiStore()
            const apiStoreData = storeToRefs(apiStore)
            const res = await request(apiStoreData.api.value.server_getSetting)
            this.serverConfig = res.data
        },
        //获取公共系统设置
        async getPublicServerConfig() {
            const apiStore = useApiStore()
            const apiStoreData = storeToRefs(apiStore)
            const res = await request(apiStoreData.staticApi.value.public_getPublicSetting)
            this.publicServerConfig = res.data
            this.acceptable_email_suffixes_arr = this.publicServerConfig.acceptable_email_suffixes.split("\n")
        },
        //修改系统设置
        async updateServerConfig(data?: object) {
            const apiStore = useApiStore()
            const apiStoreData = storeToRefs(apiStore)
            const res = await request(apiStoreData.api.value.server_updateSetting, data)
        }
    }
})