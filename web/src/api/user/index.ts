import request from '/@/utils/request';

export function useUserApi() {
    return {
        registerApi: (data: object) => {
            return request({
                url: '',
                method: 'post',
                data,
            });
        },
    };
}
