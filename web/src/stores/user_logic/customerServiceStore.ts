import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()


export const useCustomerServiceStore = defineStore('customerService', {
  state: () => ({
    customerService:{} as CustomerService,
    customerServiceList:[] as CustomerService[],
    pushCustomerServiceRequest:{} as PushCustomerServiceRequest,
  }),
  actions: {
    async getCustomerServiceList(){
      const res = await request(apiStore.userApi.getCustomerServiceList)
      this.customerServiceList = res.data
    },
    async pushCustomerService(){
      return request(apiStore.userApi.pushCustomerService,this.pushCustomerServiceRequest)
    },
    async resetSubscribeUUID(params:CustomerService){
      return request(apiStore.userApi.resetSubscribeUUID,params)
    },
    async deleteCustomerService(params:CustomerService){
      return request(apiStore.userApi.deleteCustomerService,params)
    }
  }
})