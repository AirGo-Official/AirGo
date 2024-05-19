import {defineStore} from "pinia";

const publicApiPre = "/api/public"
const adminApiPre = "/api/admin"
const userApiPre = "/api/customer"

export const useApiStore = defineStore("apiStore", {
    state: () => ({
        //固定api
        publicApi: {
            register: {path: publicApiPre + "/user/register", method: "post"} as ApiItem,
            login: {path: publicApiPre + "/user/login", method: "post"} as ApiItem,
            resetUserPassword: {path: publicApiPre + "/user/resetUserPassword", method: "post"} as ApiItem,

            getThemeConfig: {path: publicApiPre + "/server/getThemeConfig", method: "get"} as ApiItem,
            getPublicSetting: {path: publicApiPre + "/server/getPublicSetting", method: "get"} as ApiItem,

            getEmailCode: {path: publicApiPre + "/code/getEmailCode", method: "post"} as ApiItem,
            getBase64Captcha: {path: publicApiPre + "/code/getBase64Captcha", method: "get"} as ApiItem,

            getDefaultArticleList:{path: publicApiPre + "/article/getDefaultArticleList", method: "get"} as ApiItem,
        },

        adminApi:{
            // user
            newUser:{path: adminApiPre + "/user/newUser", method: "post"} as ApiItem,
            getUserList:{path: adminApiPre + "/user/getUserList", method: "post"} as ApiItem,
            updateUser:{path: adminApiPre + "/user/updateUser", method: "post"} as ApiItem,
            deleteUser:{path: adminApiPre + "/user/deleteUser", method: "delete"} as ApiItem,
            userSummary:{path: adminApiPre + "/user/userSummary", method: "post"} as ApiItem,

            // customerService
            getCustomerServiceList:{path: adminApiPre + "/customerService/getCustomerServiceList", method: "post"} as ApiItem,
            updateCustomerService:{path: adminApiPre + "/customerService/updateCustomerService", method: "post"} as ApiItem,
            deleteCustomerService:{path: adminApiPre + "/customerService/deleteCustomerService", method: "delete"} as ApiItem,

            //menu
            newMenu:{path: adminApiPre + "/menu/newMenu", method: "post"} as ApiItem,
            updateMenu:{path: adminApiPre + "/menu/updateMenu", method: "post"} as ApiItem,
            delMenu:{path: adminApiPre + "/menu/delMenu", method: "delete"} as ApiItem,
            getAllMenuList:{path: adminApiPre + "/menu/getAllMenuList", method: "get"} as ApiItem,

            // role
            newRole:{path: adminApiPre + "/role/newRole", method: "post"} as ApiItem,
            getRoleList:{path: adminApiPre + "/role/getRoleList", method: "get"} as ApiItem,
            updateRole:{path: adminApiPre + "/role/updateRole", method: "post"} as ApiItem,
            delRole:{path: adminApiPre + "/role/delRole", method: "delete"} as ApiItem,
            getAllPolicy:{path: adminApiPre + "/role/getAllPolicy", method: "get"} as ApiItem,
            getPolicyByID:{path: adminApiPre + "/role/getPolicyByID", method: "post"} as ApiItem,

            // server
            updateThemeConfig:{path: adminApiPre + "/server/updateThemeConfig", method: "post"} as ApiItem,
            getSetting:{path: adminApiPre + "/server/getSetting", method: "get"} as ApiItem,
            updateSetting:{path: adminApiPre + "/server/updateSetting", method: "post"} as ApiItem,
            getCurrentVersion:{path: adminApiPre + "/server/getCurrentVersion", method: "get"} as ApiItem,
            getLatestVersion:{path: adminApiPre + "/server/getLatestVersion", method: "get"} as ApiItem,
            updateLatestVersion:{path: adminApiPre + "/server/updateLatestVersion", method: "get"} as ApiItem,

            // node
            newNode:{path: adminApiPre + "/node/newNode", method: "post"} as ApiItem,
            getNodeList:{path: adminApiPre + "/node/getNodeList", method: "post"} as ApiItem,
            getNodeListWithTraffic:{path: adminApiPre + "/node/getNodeListWithTraffic", method: "post"} as ApiItem,
            updateNode:{path: adminApiPre + "/node/updateNode", method: "post"} as ApiItem,
            deleteNode:{path: adminApiPre + "/node/deleteNode", method: "delete"} as ApiItem,

            nodeSort:{path: adminApiPre + "/node/nodeSort", method: "post"} as ApiItem,
            createx25519:{path: adminApiPre + "/node/createx25519", method: "get"} as ApiItem,

            parseUrl:{path: adminApiPre + "/node/parseUrl", method: "post"} as ApiItem,
            newNodeShared:{path: adminApiPre + "/node/newNodeShared", method: "post"} as ApiItem,
            getNodeSharedList:{path: adminApiPre + "/node/getNodeSharedList", method: "get"} as ApiItem,
            deleteNodeShared:{path: adminApiPre + "/node/deleteNodeShared", method: "delete"} as ApiItem,
            getNodeServerStatus:{path: adminApiPre + "/node/getNodeServerStatus", method: "get"} as ApiItem,

            // shop
            newGoods:{path: adminApiPre + "/shop/newGoods", method: "post"} as ApiItem,
            getGoodsList:{path: adminApiPre + "/shop/getGoodsList", method: "get"} as ApiItem,
            updateGoods:{path: adminApiPre + "/shop/updateGoods", method: "post"} as ApiItem,
            deleteGoods:{path: adminApiPre + "/shop/deleteGoods", method: "delete"} as ApiItem,
            goodsSort:{path: adminApiPre + "/shop/goodsSort", method: "post"} as ApiItem,

            // order
            getOrderList:{path: adminApiPre + "/order/getOrderList", method: "post"} as ApiItem,
            orderSummary:{path: adminApiPre + "/order/orderSummary", method: "post"} as ApiItem,
            updateOrder:{path: adminApiPre + "/order/updateOrder", method: "post"} as ApiItem,

            // pay
            newPay:{path: adminApiPre + "/pay/newPay", method: "post"} as ApiItem,
            getPayList:{path: adminApiPre + "/pay/getPayList", method: "get"} as ApiItem,
            updatePay:{path: adminApiPre + "/pay/updatePay", method: "post"} as ApiItem,
            deletePay:{path: adminApiPre + "/pay/deletePay", method: "delete"} as ApiItem,

            // report
            getTables:{path: adminApiPre + "/report/getTables", method: "post"} as ApiItem,
            getColumn:{path: adminApiPre + "/report/getColumn", method: "post"} as ApiItem,
            reportSubmit:{path: adminApiPre + "/report/reportSubmit", method: "post"} as ApiItem,

            // article
            newArticle:{path: adminApiPre + "/article/newArticle", method: "post"} as ApiItem,
            getArticleList:{path: adminApiPre + "/article/getArticleList", method: "post"} as ApiItem,
            updateArticle:{path: adminApiPre + "/article/updateArticle", method: "post"} as ApiItem,
            deleteArticle:{path: adminApiPre + "/article/deleteArticle", method: "delete"} as ApiItem,

            // coupon
            newCoupon:{path: adminApiPre + "/coupon/newCoupon", method: "post"} as ApiItem,
            getCouponList:{path: adminApiPre + "/coupon/getCouponList", method: "post"} as ApiItem,
            updateCoupon:{path: adminApiPre + "/coupon/updateCoupon", method: "post"} as ApiItem,
            deleteCoupon:{path: adminApiPre + "/coupon/deleteCoupon", method: "delete"} as ApiItem,

            // access
            newAccess:{path: adminApiPre + "/access/newAccess", method: "post"} as ApiItem,
            getAccessList:{path: adminApiPre + "/access/getAccessList", method: "post"} as ApiItem,
            updateAccess:{path: adminApiPre + "/access/updateAccess", method: "post"} as ApiItem,
            deleteAccess:{path: adminApiPre + "/access/deleteAccess", method: "delete"} as ApiItem,

            // migration
            migrationData:{path: adminApiPre + "/migration/migrationData", method: "post"} as ApiItem,

            // ticket
            firstTicket:{path: adminApiPre + "/ticket/firstTicket", method: "post"} as ApiItem,
            getTicketList:{path: adminApiPre + "/ticket/getTicketList", method: "post"} as ApiItem,
            updateTicket:{path: adminApiPre + "/ticket/updateTicket", method: "post"} as ApiItem,
            sendTicketMessage:{path: adminApiPre + "/ticket/sendTicketMessage", method: "post"} as ApiItem,
            deleteTicket:{path: adminApiPre + "/ticket/deleteTicket", method: "delete"} as ApiItem,

        },
        userApi:{
            // user
            getUserInfo:{path: userApiPre + "/user/getUserInfo", method: "get"} as ApiItem,
            changeUserPassword:{path: userApiPre + "/user/changeUserPassword", method: "post"} as ApiItem,
            changeUserAvatar:{path: userApiPre + "/user/changeUserAvatar", method: "post"} as ApiItem,
            clockIn:{path: userApiPre + "/user/clockIn", method: "get"} as ApiItem,
            setUserNotice:{path: userApiPre + "/user/setUserNotice", method: "post"} as ApiItem,

            // customerService
            getCustomerServiceList:{path: userApiPre + "/customerService/getCustomerServiceList", method: "get"} as ApiItem,
            resetSubscribeUUID:{path: userApiPre + "/customerService/resetSubscribeUUID", method: "post"} as ApiItem,
            pushCustomerService:{path: userApiPre + "/customerService/pushCustomerService", method: "post"} as ApiItem,
            deleteCustomerService:{path: userApiPre + "/customerService/deleteCustomerService", method: "delete"} as ApiItem,

            // menu
            getMenuList:{path: userApiPre + "/menu/getMenuList", method: "get"} as ApiItem,

            // shop
            purchase:{path: userApiPre + "/shop/purchase", method: "post"} as ApiItem,
            getEnabledGoodsList:{path: userApiPre + "/shop/getEnabledGoodsList", method: "get"} as ApiItem,

            // order
            preCreateOrder:{path: userApiPre + "/order/preCreateOrder", method: "post"} as ApiItem,
            getOrderInfo:{path: userApiPre + "/order/getOrderInfo", method: "post"} as ApiItem,
            getOrderList:{path: userApiPre + "/order/getOrderList", method: "post"} as ApiItem,
            getOrderInfoWaitPay:{path: userApiPre + "/order/getOrderInfoWaitPay", method: "post"} as ApiItem,

            // pay
            getEnabledPayList:{path: userApiPre + "/pay/getEnabledPayList", method: "get"} as ApiItem,

            // article
            getArticleList:{path: userApiPre + "/article/getArticleList", method: "post"} as ApiItem,

            // ticket
            newTicket:{path: userApiPre + "/ticket/newTicket", method: "post"} as ApiItem,
            getUserTicketList:{path: userApiPre + "/ticket/getUserTicketList", method: "post"} as ApiItem,
            updateUserTicket:{path: userApiPre + "/ticket/updateUserTicket", method: "post"} as ApiItem,
            sendTicketMessage:{path: userApiPre + "/ticket/sendTicketMessage", method: "post"} as ApiItem,
            firstTicket:{path: userApiPre + "/ticket/firstTicket", method: "post"} as ApiItem,
            //traffic
            getSubTrafficList:{path: userApiPre + "/traffic/getSubTrafficList", method: "post"} as ApiItem,
            //finance
            getBalanceStatementList:{path: userApiPre + "/finance/getBalanceStatementList", method: "post"} as ApiItem,
            getCommissionStatementList:{path: userApiPre + "/finance/getCommissionStatementList", method: "post"} as ApiItem,
            getInvitationUserList:{path: userApiPre + "/finance/getInvitationUserList", method: "post"} as ApiItem,
            withdrawToBalance:{path: userApiPre + "/finance/withdrawToBalance", method: "get"} as ApiItem,
            getCommissionSummary:{path: userApiPre + "/finance/getCommissionSummary", method: "get"} as ApiItem,
        },
    }),
    actions: {

    }

})