declare interface SysUser {
    created_at: string;
    updated_at: string;
    id: number;
    uuid: number;
    user_name: string;
    nick_name: string;
    password: string;
    avatar: string;
    enable: boolean;
    invitation_code: string;
    referrer_user_id: number;
    balance: number;
    tg_id: number
    role_group: RoleInfo[];	//角色组
    orders: [];      //订单

    //通知参数
    enable_tg_bot:boolean
    enable_email:boolean
    enable_web_mail:boolean
    when_purchased:boolean
    when_service_almost_expired:boolean
    when_balance_changed:boolean
}

// 登录
declare interface LoginForm {
    user_name: string
    password: string
    email_code: string
    base64_captcha: Base64CaptchaInfo
}

// 注册
declare interface RegisterForm {
    user_name: string
    email_suffix: string
    password: string
    re_password: string
    email_code: string
    referrer_code: string
    base64_captcha: Base64CaptchaInfo
}
//
declare interface Base64CaptchaInfo {
    id: string
    b64s: string

}
declare interface UserSummary {
    date:string
    register_total:number
}