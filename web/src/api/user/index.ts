import request from '/@/utils/request';

/**
 * （不建议写成 request.post(xxx)，因为这样 post 时，无法 params 与 data 同时传参）
 *
 * 登录api接口集合
 * @method signIn 用户登录
 * @method signOut 用户退出登录
 */
export function useUserApi() {
    return {
        registerApi: (data: object) => {
            return request({
                url: '/user/register',
                method: 'post',
                data,
            });
        },
        signIn: (data: object) => {
            return request({
                url: '/user/login',
                method: 'post',
                data,
            });
        },
        signOut: (data: object) => {
            return request({
                url: '/user/signOut',
                method: 'post',
                data,
            });
        },
        //修改混淆
        changeHostApi: (data?: object) => {
            return request({
                url: '/user/changeSubHost',
                method: 'post',
                data,
            });
        },
        //获取自身信息
        getUserInfoApi: () => {
            return request({
                url: '/user/getUserInfo',
                method: 'get',
            });
        },
        //获取用户列表
        getUserListApi: (data?: object) => {
            return request({
                url: '/user/getUserList',
                method: 'post',
                data: data
            });
        },
        //新建用户
        newUserApi: (data?: object) => {
            return request({
                url: '/user/newUser',
                method: 'post',
                data
            });
        },
        //修改用户
        updateUserApi: (data?: object) => {
            return request({
                url: '/user/updateUser',
                method: 'post',
                data
            });
        },
        //删除用户
        deleteUserApi: (data?: object) => {
            return request({
                url: '/user/deleteUser',
                method: 'post',
                data
            });
        },
        //查询用户
        findUserApi: (data?: object) => {
            return request({
                url: '/user/findUser',
                method: 'post',
                data
            });
        },
        //修改密码
        changePasswordApi: (data?: object) => {
            return request({
                url: '/user/changeUserPassword',
                method: 'post',
                data
            });
        },
        //确认重置密码
        resetPasswordApi: (data?: object) => {
            return request({
                url: '/user/resetUserPassword',
                method: 'post',
                data
            });
        },
        //重置订阅
        resetSubApi: () => {
            return request({
                url: '/user/resetSub',
                method: 'get',
            });
        },
    };
}
