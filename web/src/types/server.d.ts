declare interface Server {
    created_at: string;
    updated_at: string;
    id: number;

    subscribe: Subscribe;
    email: Email;
    security: Security;
    notice: Notice;
}
declare interface Notice {
    bot_token: string;
    tg_admin: string;
    tg_socks5 :string;
    when_user_registered: boolean;
    when_user_purchased: boolean;
    when_node_offline: boolean;
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

declare interface Subscribe {
    enable_register: boolean;
    enable_email_code: boolean;
    enable_login_email_code: boolean;
    acceptable_email_suffixes: string;
    is_multipoint: boolean;

    backend_url: string;
    frontend_url: string;
    api_prefix: string;

    sub_name: string;
    tek: string;
    default_goods: number;
    enabled_rebate: boolean;    //是否开启返利
    rebate_rate: number;        //返利率
    enabled_deduction: boolean; //是否开启旧套餐抵扣
    deduction_threshold: number;//旧套餐抵扣阈值,大于该值则抵扣

    enabled_clock_in: boolean
    clock_in_min_traffic: number
    clock_in_max_traffic: number
    clock_in_min_day: number
    clock_in_max_day: number
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