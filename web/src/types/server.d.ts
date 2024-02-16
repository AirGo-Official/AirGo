declare interface Server {
    created_at: string;
    updated_at: string;
    id: number;

    website: Website;
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

declare interface Website {
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
    enabled_clock_in: boolean
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