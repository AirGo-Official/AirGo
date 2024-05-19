declare interface Server {
    created_at: string;
    updated_at: string;
    id: number;

    website: Website;
    subscribe: Subscribe
    email: Email;
    security: Security;
    notice: Notice;
    finance: Finance
}
declare interface Subscribe{
    backend_url: string;
    sub_name: string;
    tek: string;
    subscribe_domain_bind_request:boolean
    surge_rule:string
    clash_rule:string
}
declare interface Notice {
    enable_tg_bot:boolean
    enable_email:boolean
    enable_web_mail:boolean
    admin_id:string
    bot_token:string
    tg_socks5:string
    when_user_registered:boolean
    when_user_purchased:boolean
    when_node_offline:boolean
    when_new_ticket:boolean
}
declare interface Security {
    captcha: Captcha
    jwt: Jwt
    rate_limit_params: RateLimitParams
}

declare interface Jwt {
    signing_key: string;
    expires_time: string;
    buffer_time: string;
    issuer: string;
}

declare interface Website {
    enable_register: boolean;
    enable_base64_captcha:boolean; //是否开启注册图片验证码
    enable_email_code: boolean; //是否开启注册email 验证码
    enable_login_email_code: boolean; //是否开启登录email 验证码
    acceptable_email_suffixes: string;
    is_multipoint: boolean;
    frontend_url: string;
    enable_swagger_api:boolean
    enable_assets_api:boolean
}

declare interface Captcha {
    key_long: number;
    img_width: number;
    img_height: number;
    open_captcha: number;
    open_captcha_time_out: number;
}
declare interface Email {
    email_from: string;
    email_from_alias: string;
    email_secret: string;
    email_host: string;
    email_port: number;
    email_is_ssl: boolean;
    email_nickname: string;
    email_subject: string;
    email_content: string;
}

declare interface RateLimitParams {
    ip_role_param: number;
    visit_param: number;

}
declare interface Finance {
    enable_invitation_commission:boolean
    commission_rate:number //佣金率, 范围 0~1, 佣金 = 订单金额 * 佣金率 ( 100.50 * 0.50 )
    withdraw_threshold:number //提取到余额的阈值
    enable_lottery:boolean
    jackpot:JackpotItem[]
}

declare interface JackpotItem {
    balance:number
    weight:number
}