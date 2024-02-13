import {defineStore, storeToRefs} from "pinia";

import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

export const useAdminServerStore = defineStore("serverAdminStore", {
    state: () => ({
        serverConfig: {
            created_at: '',
            updated_at: '',
            id: 0,
            security: {
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
                acceptable_email_suffixes: '',
                is_multipoint: true,
                sub_name: '',
                backend_url: '',
                frontend_url: '',
                api_prefix: '',
                tek: '',
                default_goods: 0,
                enabled_rebate: true,    //是否开启返利
                rebate_rate: 0,          //返利率
                enabled_deduction: true, //是否开启旧套餐抵扣
                deduction_threshold: 0,  //旧套餐抵扣阈值,大于该值则抵扣
                enabled_clock_in: true,
                clock_in_min_traffic: 100,
                clock_in_max_traffic: 1000,
                clock_in_min_day: 0,
                clock_in_max_day: 0,
            } as Subscribe,

            email: {
                email_from: '',
                email_from_alias: 'Admin',
                email_secret: '',
                email_host: '',
                email_port: 465,
                email_is_ssl: false,
                email_nickname: '管理员',
                email_subject: '',
                email_content: '',
            } as Email,
            notice: {
                bot_token: '',
                tg_admin: '',
                tg_socks5: '',
                when_user_registered: false,
                when_user_purchased: false,
                when_node_offline: false,
            } as Notice,

        } as Server,
    }),
    actions: {
        //获取系统设置
        async getServerConfig() {
            const res = await request(apiStoreData.adminApi.value.getSetting)
            this.serverConfig = res.data
        },
        //修改系统设置
        async updateServerConfig(params?: object) {
            return  request(apiStoreData.adminApi.value.updateSetting, params)
        }
    }
})