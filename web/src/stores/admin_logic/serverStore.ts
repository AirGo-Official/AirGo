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
            website: {
                enable_register: true,
                enable_email_code: false,
                enable_login_email_code: false,
                acceptable_email_suffixes: '',
                is_multipoint: true,
                frontend_url: '',
                api_prefix: '',
                enabled_clock_in: true,
            } as Website,
            subscribe:{
                backend_url: '',
                sub_name: '',
                tek: '',
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
                enable_tg_bot:false,
                enable_email:false,
                enable_web_mail:false,
                admin_id:'',
                bot_token:'',
                tg_socks5:'',
                when_user_registered:false,
                when_user_purchased:false,
                when_node_offline:false,
            } as Notice,

        } as Server,
        version:{
            currentVersion:{
                version:''
            },
            latestVersion:{
                version:''
            }
        }
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
        },
        //获取AirGo核心当前版本
        async getCurrentVersion(){
            const res = await request(apiStoreData.adminApi.value.getCurrentVersion)
            this.version.currentVersion = res.data
        },
        //获取AirGo核心最新版本
        async getLatestVersion(){
            const res = await request(apiStoreData.adminApi.value.getLatestVersion)
            this.version.latestVersion = res.data
        },
        //升级AirGo核心最新版本
        async updateLatestVersion(){
            const res = await request(apiStoreData.adminApi.value.updateLatestVersion)
            // TODO 应提醒用户正在升级中
        }
    }
})