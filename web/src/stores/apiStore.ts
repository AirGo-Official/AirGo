import {defineStore} from "pinia";

const apiPre = "/api"

export const useApiStore = defineStore("apiStore", {
    state: () => ({
        //固定api
        staticApi: {
            public_getEmailCode: {path: apiPre + "/public/getEmailCode", method: "post"} as ApiItem,
            public_getBase64Captcha: {path: apiPre + "/public/getBase64Captcha", method: "get"} as ApiItem,
            public_epayNotify: {path: apiPre + "/public/epayNotify", method: "get"} as ApiItem,
            public_alipayNotify: {path: apiPre + "/public/alipayNotify", method: "post"} as ApiItem,
            public_queryPackage: {path: apiPre + "/public/queryPackage", method: "get"} as ApiItem,

            public_getThemeConfig: {path: apiPre + "/public/getThemeConfig", method: "get"} as ApiItem,
            public_getPublicSetting: {path: apiPre + "/public/getPublicSetting", method: "get"} as ApiItem,

            user_register: {path: apiPre + "/user/register", method: "post"} as ApiItem,
            user_login: {path: apiPre + "/user/login", method: "post"} as ApiItem,
            user_getSub: {path: apiPre + "/user/getSub", method: "get"} as ApiItem,
            user_resetUserPassword: {path: apiPre + "/user/resetUserPassword", method: "post"} as ApiItem,
        },
        //api
        api: {
            //user
            user_changeSubHost: {path: apiPre + "/user/changeSubHost", method: "post"} as ApiItem,
            user_getUserInfo: {path: apiPre + "/user/getUserInfo", method: "get"} as ApiItem,
            user_changeUserPassword: {path: apiPre + "/user/changeUserPassword", method: "post"} as ApiItem,
            user_resetSub: {path: apiPre + "/user/resetSub", method: "get"} as ApiItem,

            user_getUserList: {path: apiPre + "/user/getUserList", method: "post"} as ApiItem,
            user_newUser: {path: apiPre + "/user/newUser", method: "post"} as ApiItem,
            user_updateUser: {path: apiPre + "/user/updateUser", method: "post"} as ApiItem,
            user_deleteUser: {path: apiPre + "/user/deleteUser", method: "post"} as ApiItem,

            //menu
            menu_getRouteList: {path: apiPre + "/menu/getRouteList", method: "get"} as ApiItem,
            menu_getRouteTree: {path: apiPre + "/menu/getRouteTree", method: "get"} as ApiItem,
            menu_getAllRouteList: {path: apiPre + "/menu/getAllRouteList", method: "get"} as ApiItem,
            menu_getAllRouteTree: {path: apiPre + "/menu/getAllRouteTree", method: "get"} as ApiItem,
            menu_newDynamicRoute: {path: apiPre + "/menu/newDynamicRoute", method: "post"} as ApiItem,
            menu_delDynamicRoute: {path: apiPre + "/menu/delDynamicRoute", method: "post"} as ApiItem,
            menu_updateDynamicRoute: {path: apiPre + "/menu/updateDynamicRoute", method: "post"} as ApiItem,
            menu_findDynamicRoute: {path: apiPre + "/menu/findDynamicRoute", method: "post"} as ApiItem,

            //role
            role_getRoleList: {path: apiPre + "/role/getRoleList", method: "post"} as ApiItem,
            role_modifyRoleInfo: {path: apiPre + "/role/modifyRoleInfo", method: "post"} as ApiItem,
            role_addRole: {path: apiPre + "/role/addRole", method: "post"} as ApiItem,
            role_delRole: {path: apiPre + "/role/delRole", method: "post"} as ApiItem,
            //系统设置
            system_updateThemeConfig: {path: apiPre + "/system/updateThemeConfig", method: "post"} as ApiItem,
            system_getSetting: {path: apiPre + "/system/getSetting", method: "get"} as ApiItem,
            system_updateSetting: {path: apiPre + "/system/updateSetting", method: "post"} as ApiItem,
            system_createx25519: {path: apiPre + "/system/createx25519", method: "get"} as ApiItem,

            //节点
            node_getAllNode: {path: apiPre + "/node/getAllNode", method: "get"} as ApiItem,
            node_newNode: {path: apiPre + "/node/newNode", method: "post"} as ApiItem,
            node_deleteNode: {path: apiPre + "/node/deleteNode", method: "post"} as ApiItem,
            node_updateNode: {path: apiPre + "/node/updateNode", method: "post"} as ApiItem,
            node_getTraffic: {path: apiPre + "/node/getTraffic", method: "post"} as ApiItem,
            node_nodeSort: {path: apiPre + "/node/nodeSort", method: "post"} as ApiItem,
            node_newNodeShared: {path: apiPre + "/node/newNodeShared", method: "post"} as ApiItem,
            node_getNodeSharedList: {path: apiPre + "/node/getNodeSharedList", method: "get"} as ApiItem,
            node_deleteNodeShared: {path: apiPre + "/node/deleteNodeShared", method: "post"} as ApiItem,

            //shop
            shop_preCreatePay: {path: apiPre + "/shop/preCreatePay", method: "post"} as ApiItem,
            shop_purchase: {path: apiPre + "/shop/purchase", method: "post"} as ApiItem,
            shop_getAllEnabledGoods: {path: apiPre + "/shop/getAllEnabledGoods", method: "get"} as ApiItem,
            shop_getAllGoods: {path: apiPre + "/shop/getAllGoods", method: "get"} as ApiItem,
            shop_newGoods: {path: apiPre + "/shop/newGoods", method: "post"} as ApiItem,
            shop_deleteGoods: {path: apiPre + "/shop/deleteGoods", method: "post"} as ApiItem,
            shop_updateGoods: {path: apiPre + "/shop/updateGoods", method: "post"} as ApiItem,
            shop_goodsSort: {path: apiPre + "/shop/goodsSort", method: "post"} as ApiItem,

            //order
            order_getOrderInfo: {path: apiPre + "/order/getOrderInfo", method: "post"} as ApiItem,
            order_getOrderByUserID: {path: apiPre + "/order/getOrderByUserID", method: "post"} as ApiItem,
            order_getAllOrder: {path: apiPre + "/order/getAllOrder", method: "post"} as ApiItem,
            order_completedOrder: {path: apiPre + "/order/completedOrder", method: "post"} as ApiItem,
            order_getMonthOrderStatistics: {path: apiPre + "/order/getMonthOrderStatistics", method: "post"} as ApiItem,

            //pay
            pay_getEnabledPayList: {path: apiPre + "/pay/getEnabledPayList", method: "get"} as ApiItem,
            pay_getPayList: {path: apiPre + "/pay/getPayList", method: "get"} as ApiItem,

            pay_newPay: {path: apiPre + "/pay/newPay", method: "post"} as ApiItem,
            pay_deletePay: {path: apiPre + "/pay/deletePay", method: "post"} as ApiItem,
            pay_updatePay: {path: apiPre + "/pay/updatePay", method: "post"} as ApiItem,

            //casbin
            casbin_getAllPolicy: {path: apiPre + "/casbin/getAllPolicy", method: "get"} as ApiItem,
            casbin_getPolicyByRoleIds: {path: apiPre + "/casbin/getPolicyByRoleIds", method: "post"} as ApiItem,
            casbin_updateCasbinPolicy: {path: apiPre + "/casbin/updateCasbinPolicy", method: "post"} as ApiItem,

            //websocket
            websocket_msg: {path: apiPre + "/websocket/msg", method: "get"} as ApiItem,

            //upload
            upload_newPictureUrl: {path: apiPre + "/upload/newPictureUrl", method: "post"} as ApiItem,
            upload_getPictureList: {path: apiPre + "/upload/getPictureList", method: "post"} as ApiItem,

            //report
            report_getDB: {path: apiPre + "/report/getDB", method: "get"} as ApiItem,
            report_getTables: {path: apiPre + "/report/getTables", method: "post"} as ApiItem,
            report_getColumn: {path: apiPre + "/report/getColumn", method: "post"} as ApiItem,
            report_reportSubmit: {path: apiPre + "/report/reportSubmit", method: "post"} as ApiItem,

            //article
            article_newArticle: {path: apiPre + "/article/newArticle", method: "post"} as ApiItem,
            article_deleteArticle: {path: apiPre + "/article/deleteArticle", method: "post"} as ApiItem,
            article_updateArticle: {path: apiPre + "/article/updateArticle", method: "post"} as ApiItem,
            article_getArticle: {path: apiPre + "/article/getArticle", method: "post"} as ApiItem,

            //coupon
            coupon_newCoupon: {path: apiPre + "/coupon/newCoupon", method: "post"} as ApiItem,
            coupon_deleteCoupon: {path: apiPre + "/coupon/deleteCoupon", method: "post"} as ApiItem,
            coupon_updateCoupon: {path: apiPre + "/coupon/updateCoupon", method: "post"} as ApiItem,
            coupon_getCoupon: {path: apiPre + "/coupon/getCoupon", method: "post"} as ApiItem,

            //isp
            isp_sendCode: {path: apiPre + "/isp/sendCode", method: "post"} as ApiItem,
            isp_ispLogin: {path: apiPre + "/isp/ispLogin", method: "post"} as ApiItem,
            isp_getMonitorByUserID: {path: apiPre + "/isp/getMonitorByUserID", method: "post"} as ApiItem,


        },

    }),
    actions: {}

})