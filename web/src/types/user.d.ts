declare interface SysUser {
    created_at: string;
    updated_at: string;
    id: number;
    uuid: number;
    user_name: string;
    nick_name: string;
    password: string;
    avatar: string;
    phone: string;
    email: string;
    enable: boolean;
    invitation_code: string;
    referrer_code: string;
    remain: number;
    role_group: RowRoleType[];	//角色组
    orders: [];      //订单
    subscribe_info: { //附加订阅信息
        host: string;
        client_ip: string;
        sub_status: boolean;
        subscribe_url: string;
        goods_id: int;
        goods_subject: string;
        expired_at: string;
        t: number;
        u: number;
        d: number;
        reset_day: number;
        node_speedlimit: number;
        node_connector: number;
    }
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

declare interface Base64CaptchaInfo{
        id:string;
        b64s:string;
}
