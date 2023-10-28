declare interface Isp {
    created_at: string;
    updated_at: string;
    id: number;
    user_id: number;
    mobile: string;
    isp_type: string;
    status: boolean;
    unicom_config: UnicomConfig;
    telecom_config:TelecomConfig;
}

declare interface UnicomConfig {
    version: string;
    app_id: string;
    cookie: string;
    unicomMobile: string;
    password: string;
}

declare interface TelecomConfig {
    phoneNum: string;
    telecomPassword: string;

    timestamp: string;
    loginAuthCipherAsymmertric: string;

    deviceUid: string;
    telecomToken: string;

}