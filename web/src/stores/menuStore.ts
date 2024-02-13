//路由-store
import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

export const useMenuStore = defineStore("menuStore", {
    state: () => ({
        //当前角色的动态路由list
        routesListSate: {
            routesList: [],
            isColumnsMenuHover: false,
            isColumnsNavHover: false,
        } as RoutesListState,
    }),
    actions: {
        // 设置动态路由到 pinia routesList 中
        async setRoutesList(data: Array<Route>) {
            this.routesListSate.routesList = data;
        },
        async setColumnsMenuHover(bool: Boolean) {
            this.routesListSate.isColumnsMenuHover = bool;
        },
        async setColumnsNavHover(bool: Boolean) {
            this.routesListSate.isColumnsNavHover = bool;
        },
    }
})