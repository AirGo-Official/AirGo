import {defineStore, storeToRefs} from 'pinia';
import {Local, Session} from '/@/utils/storage';
import {usePublicStore} from "/@/stores/publicStore";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
const publicStore = usePublicStore()

export const useUserStore = defineStore('userStore', {
    state: () => ({
        //登录页面数据
        loginData: {
            user_name: '',
            password: '',
            email_code: '',
            base64_captcha: {
                id: '',
                b64s: '',
            }
        } as LoginForm,
        //注册数据
        registerData: {
            user_name: '',
            email_suffix: '@qq.com',
            password: '',
            re_password: '',
            email_code: '',
            referrer_code: '',
            base64_captcha: {
                id: '',
                b64s: '',
            }
        } as RegisterForm,
        //修改用户头像
        userAvatar:{
            avatar:'',
        },
        //全局用户信息
        userInfos: {
            created_at: '',
            updated_at: '',
            id: 0,
            uuid: 0,
            user_name: '',
            nick_name: '',
            password: '',
            avatar: '',
            enable: true,
            invitation_code: '',
            referrer_user_id: 0,
            balance: 0,
            role_group: [] as RoleInfo[],	//角色组
            orders: [],      //订单
            //通知参数
            enable_tg_bot:false,
            enable_email:false,
            enable_web_mail:false,
            when_purchased:false,
            when_service_almost_expired:false,
            when_balance_changed:false,
        } as SysUser,
    }),
    actions: {
        //注册
        async register(form?: object) {
            if (Session.get('invitation')) {
                this.registerData.referrer_code = Session.get('invitation')
            }
            this.registerData.base64_captcha.id = publicStore.base64CaptchaData.id
            return  request(apiStore.publicApi.register, this.registerData)
        },
        //登录
        async userLogin(params?: any) {
            const res: any = await request(apiStore.publicApi.login, params)
            //保存用户信息到pinia
            this.userInfos = res.data.user;
            //保存用户信息到Session
            Session.set("userInfos", res.data.user)
            //保存token
            Local.set("token", res.data.token)
        },
        //获取自身信息
        async getUserInfo() {
            const res = await request(apiStore.userApi.getUserInfo)
            this.userInfos = res.data
            Session.set("userInfos", res.data)
        },
        //修改密码
        async changePassword() {
            return  request(apiStore.userApi.changeUserPassword, this.registerData)
        },
        //修改头像
        async changeAvatar() {
            return  request(apiStore.userApi.changeUserAvatar, this.userAvatar)
        },
        //重置密码
        async submitResetPassword() {
            return request(apiStore.publicApi.resetUserPassword, this.loginData)
        },
        async setUserNotice() {
            return  request(apiStore.userApi.setUserNotice, this.userInfos)
        },
        async clockIn(){
            return  await request(apiStore.userApi.clockIn)
        }
    },
});
