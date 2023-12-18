import {defineStore} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
export const useTrafficStore = defineStore("trafficStore", {
        state: () => ({
            userTrafficLog: {
                data: [] as UserTrafficLog[],
                total: 0,
            },
            currentUserTrafficLog:{} as UserTrafficLog,

        }),
    getters:{
        //上行流量，GB
        upTraffic: (state): number => {
            return +(state.currentUserTrafficLog.u / 1024 / 1024 / 1024).toFixed(2)
        },
        //下行流量，GB
        downTraffic: (state): number => {
            return +(state.currentUserTrafficLog.d / 1024 / 1024 / 1024).toFixed(2)
        },

    },
        actions: {
            async getUserTrafficLog(params: object) {
                const res = await request(apiStore.api.user_getUserTraffic, params)
                this.currentUserTrafficLog = res.data
            },
            async getAllUserTrafficLog(params: object) {
                const res = await request(apiStore.api.user_getAllUserTraffic, params)
                this.userTrafficLog = res.data
            }
        }
    }
)