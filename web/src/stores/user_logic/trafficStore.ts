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
    async getSubTrafficList(){
      const res = await request(apiStore.userApi.getSubTrafficList)
      const list:UserTrafficLog[] = res.data
      if (list.length > 0){
       this.trafficHandler(list)
      }
      return res
      // this.subTrafficList=res.data
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