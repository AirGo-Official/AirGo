import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useAccessStore = defineStore("accessStore", {
    state: () => ({
        params: {
            table_name: 'access',
            field_params_list: [
                {field: 'name', condition: 'like', condition_value: '',},
            ] as FieldTableNew[],
            pagination: {
                page_num: 1,
                page_size: 30,
            } as Pagination,//分页参数
        } as FieldParams,
        routes_list: {
            total: 0,
            data: [] as Access[]
        },
        current_routes: {} as Access,

    }),
    actions: {
        async getRoutesList(data?:object) {
            const res = await request(apiStoreData.api.value.access_getRoutesList,data)
            this.routes_list = res.data
        },
        async newRoutes(data: object) {
            const res = await request(apiStoreData.api.value.access_newRoutes, data)

        },
        async deleteRoutes(data: object) {
            const res = await request(apiStoreData.api.value.access_deleteRoutes, data)

        },

        async updateRoutes(data: object) {
            const res = await request(apiStoreData.api.value.access_updateRoutes, data)

        }
    }
})