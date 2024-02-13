import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()


export const useAdminCustomerServiceStore = defineStore('adminCustomerServiceStore', {
  state: () => ({
    customerServiceList:{
      total:0,
      data:[] as CustomerService[],
    },
    currentCustomerService:{} as CustomerService,

  }),
  actions: {
    //
    async getCustomerServiceList(params:object){
     const res = await request(apiStore.adminApi.getCustomerServiceList,params)
      this.customerServiceList.data = res.data
    },
    //
    async updateCustomerService(params:object){
      return request(apiStore.adminApi.updateCustomerService,params)
    },
  }
})