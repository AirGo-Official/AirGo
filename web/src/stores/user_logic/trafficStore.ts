import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { DateStrHandler } from "/@/utils/formatTime";
const apiStore = useApiStore()


export const useTrafficStore = defineStore('trafficStore', {
  state: () => ({
    subTrafficList:[] as UserTrafficLog[],
    trafficLineChart:{
      u:[] as number[],
      d:[] as number[],
      xAxis:[] as string[],
    },
  }),
  actions: {
    async getSubTrafficList(params:UserTrafficLog){
      this.subTrafficList = [] as UserTrafficLog[]//清除上一个查询的数据
      const res = await request(apiStore.userApi.getSubTrafficList,params)
      if ( res.data === null){
        return
      }
      const list:UserTrafficLog[] = res.data
      if (list.length > 0){
       this.trafficHandler(list)
      }
    },
    trafficHandler(list:UserTrafficLog[]){
      this.trafficLineChart.u = []
      this.trafficLineChart.d = []
      this.trafficLineChart.xAxis=[]
      list.forEach((value: UserTrafficLog, index: number, array: UserTrafficLog[])=>{
        this.trafficLineChart.xAxis.push(DateStrHandler(value.created_at))
        this.trafficLineChart.u.push(Number((value.u/1024/1024/1024).toFixed(2)))
        this.trafficLineChart.d.push(Number((value.d/1024/1024/1024).toFixed(2)))
      })
    }
  }
})