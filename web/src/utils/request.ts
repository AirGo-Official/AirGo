import axios, {AxiosInstance, AxiosResponse} from 'axios';
import {ElMessage, ElMessageBox} from 'element-plus';
import {Local} from '/@/utils/storage';
import qs from 'qs';

let apiUrl = import.meta.env.VITE_API_URL

if (apiUrl === '') {
    // console.log("apiUrl===''")
    apiUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port
}

// 配置新建一个 axios 实例
const service: AxiosInstance = axios.create({
    baseURL: apiUrl,
    timeout: 30000,
    headers: {'Content-Type': 'application/json'},
    // headers: {"Content-Type": "application/x-www-form-urlencoded"},
    paramsSerializer: {
        serialize(params) {
            return JSON.stringify(params);//json提交
        },
    },
});

// 请求拦截器
service.interceptors.request.use(
    (config) => {
        // console.log("请求拦截器 config:", config)
        if (Local.get('token')) {
            config.headers!['Authorization'] = `${Local.get('token')}`;
        }
        if (config.method === "get") { //get等URL传参，重写请求头和序列化
            config.headers['Content-Type'] = "application/x-www-form-urlencoded"
            config.paramsSerializer = {
                serialize(params) {
                    return qs.stringify(params, {allowDots: true});
                }
            }
        }
        return config;
    },
    (error) => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);

// 响应拦截器
service.interceptors.response.use(
    (response) => {
        // 对响应数据做点什么
        const res = response.data;
        console.log("响应数据：", res);
        if (res.code && res.code !== 0) {
            // `token` 过期或者账号已在别处登录
            if (res.code === 401) {
                // Session.clear(); // 清除浏览器全部临时缓存
                Local.clear()
                window.location.href = '/'; // 去登录页
                ElMessageBox.alert('请重新登录', '提示', {})
                    .then(() => {
                    })
                    .catch(() => {
                    });
            }
            ElMessage.error(res.msg);
            return Promise.reject(service.interceptors.response);
        } else {
            return res;
        }
    },
    (error) => {
        //console.log("响应错误")
        // 对响应错误做点什么
        if (error.message.indexOf('timeout') != -1) {
            ElMessage.error('网络超时');
        } else if (error.message == 'Network Error') {
            ElMessage.error('网络连接错误');
        } else {
            if (error.response.data) ElMessage.error(error.response.statusText);
            else ElMessage.error('接口路径找不到');
        }
        return Promise.reject(error);
    }
);


// 导出 axios 实例
export default service;

export function request(apiItem: ApiItem, data?: any):Promise<AxiosResponse<any>> {
    if (apiItem.method === "get" || apiItem.method === "GET"){
        return service({
                        url: apiItem.path,
                        method: apiItem.method,
                        params: data,
                    })
    }
    return service({
        url: apiItem.path,
        method: apiItem.method,
        data: data,
    })
}




