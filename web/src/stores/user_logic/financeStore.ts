import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { Session } from "/@/utils/storage";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

export const useFinanceStore = defineStore("financeStore", {
  state: () => ({
    commissionSummary:{
      total_invitation:0,
      total_commission_amount:'',
      pending_withdrawal_amount:'',
      total_consumption_amount:'',
    },
    balanceStatementList:{
      total:0,
      data:[] as BalanceStatement[],

    },
    commissionStatementList:{
      total:0,
      data:[] as CommissionStatement[],
    },
    invitationUserList:{
      total:0,
      data: [] as SysUser[],
    }

  }),
  actions: {
    async getBalanceStatementList(params:QueryParams){
      const res = await request(apiStore.userApi.getBalanceStatementList,params)
      this.balanceStatementList = res.data

    },
    async getCommissionStatementList(params:QueryParams){
      const res = await request(apiStore.userApi.getCommissionStatementList,params)
      this.commissionStatementList = res.data
    },
    async getInvitationUserList(params:QueryParams){
      const res = await request(apiStore.userApi.getInvitationUserList,params)
      this.invitationUserList = res.data
    },
    async getCommissionSummary(){
      const res = await request(apiStore.userApi.getCommissionSummary)
      this.commissionSummary = res.data
    },
    async withdrawToBalance(){
      return request(apiStore.userApi.withdrawToBalance)
    }

  }

})